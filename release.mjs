#!/usr/bin/env zx

/*
* Script to release the seats.io go lib.
*   - changes the version number in README.md
*   - changes the version number in build.gradle
*   - creates the release in Gihub (using gh cli)
*
*
* Prerequisites:
*   - zx installed (https://github.com/google/zx)
*   - gh cli installed (https://cli.github.com/)
*
* Usage:
*   yarn zx ./release.mjs -v major/minor -n "release notes"
* */

// don't output the commands themselves
$.verbose = true

import { resolve } from 'path'
import { readdir } from 'fs/promises'

const semver = require('semver')

const versionToBump = getVersionToBump()
const latestReleaseTag = await fetchLatestReleasedVersionNumber()
const latestVersion = removeLeadingV(latestReleaseTag)
const nextVersion = await determineNextVersionNumber(latestVersion)
const nextMajorVersion = await getMajorVersion(nextVersion)

await assertChangesSinceRelease(latestReleaseTag)
await bumpVersionInFiles()
await commitAndPush()
await release()
await updateGoPackageRepo()

function getVersionToBump() {
    if (!argv.v || !(argv.v === 'minor' || argv.v === 'major')) {
        throw new Error ("Please specify -v major/minor")
    }
    return argv.v
}

function removeLeadingV(tagName) {
    if (tagName.startsWith('v')) {
        return tagName.substring(1)
    }
    return tagName
}

async function fetchLatestReleasedVersionNumber() {
    let result = await $`gh release view --json tagName`
    return JSON.parse(result).tagName
}

async function determineNextVersionNumber(previous) {
    return semver.inc(previous, versionToBump)
}

async function getMajorVersion(fullVersion) {
    return semver.major(fullVersion)
}

async function bumpVersionInFiles() {
    const currentMajorVersion = await getMajorVersion(latestVersion)

    await replaceInFile("README.md", `github.com/seatsio/seatsio-go/v${currentMajorVersion} v${latestVersion}`, `github.com/seatsio/seatsio-go/v${nextMajorVersion} v${nextVersion}`)
    if (nextMajorVersion > currentMajorVersion) {
        await replaceInFile("go.mod", `module github.com/seatsio/seatsio-go/v${currentMajorVersion}`, `module github.com/seatsio/seatsio-go/v${nextMajorVersion}`)
        await replaceInFile("README.md", `(https://pkg.go.dev/github.com/seatsio/seatsio-go/v${currentMajorVersion})`, `(https://pkg.go.dev/github.com/seatsio/seatsio-go/v${nextMajorVersion})`)
        await replaceInFile("README.md", `"github.com/seatsio/seatsio-go/v${currentMajorVersion}`, `"github.com/seatsio/seatsio-go/v${nextMajorVersion}`)

        for await (const filePath of getFiles('.')) {
            if (filePath.endsWith(".go")) {
                await replaceInFile(filePath, `"github.com/seatsio/seatsio-go/v${currentMajorVersion}`, `"github.com/seatsio/seatsio-go/v${nextMajorVersion}`, false)
            }
        }
    }
}

async function* getFiles(dir) {
    const dirents = await readdir(dir, { withFileTypes: true });
    for (const dirent of dirents) {
        const res = resolve(dir, dirent.name);
        if (dirent.isDirectory()) {
            yield* getFiles(res);
        } else {
            yield res;
        }
    }
}

async function replaceInFile(filename, latestVersion, nextVersion, validate = true) {
    return await fs.readFile(filename, 'utf8')
        .then(text => {
            if (validate && text.indexOf(latestVersion) < 0) {
                throw new Error('Not the correct version. Could not find ' + latestVersion + ' in ' + filename)
            }
            return text
        })
        .then(text => text.replaceAll(latestVersion, nextVersion))
        .then(text => fs.writeFileSync(filename, text))
        .then(() => gitAdd(filename))
}

async function gitAdd(filename) {
    return await $`git add ${filename}`
}

async function commitAndPush() {
    await $`git commit -m "version bump"`
    await $`git push origin main`
}

async function getCurrentCommitHash() {
    return (await $`git rev-parse HEAD`).stdout.trim()
}

async function getCommitHashOfTag(tag) {
    return (await $`git rev-list -n 1 ${tag}`).stdout.trim()
}

async function assertChangesSinceRelease(releaseTag) {
    let mainCommitHash = await getCurrentCommitHash()
    let releaseCommitHash = await getCommitHashOfTag(releaseTag)
    if(mainCommitHash === releaseCommitHash) {
        throw new Error("No changes on main since release tagged " + releaseTag)
    }
}

async function release() {
    const newTag = 'v' + nextVersion
    return await $`gh release create ${newTag} --generate-notes`.catch(error => {
        console.error('something went wrong while creating the release. Please revert the version change!')
        throw error
    })
}

async function updateGoPackageRepo() {
    const newTag = 'v' + nextVersion
    return await $`curl -X POST https://pkg.go.dev/fetch/github.com/seatsio/seatsio-go/v${nextMajorVersion}@${newTag}`.catch(error => {
        console.error('pkg.go.dev could not be updated')
        throw error
    })
}
