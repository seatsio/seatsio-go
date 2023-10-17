#!/usr/bin/env zx

/*
* Script to release the seats.io java lib.
*   - changes the version number in README.md
*   - changes the version number in build.gradle
*   - creates the release in Gihub (using gh cli)
*
*
* Prerequisites:
*   - zx installed (https://github.com/google/zx)
*   - gh cli installed (https://cli.github.com/)
*   - semver cli installed (https://github.com/fsaintjacques/semver-tool)
*
* Usage:
*   zx ./release.mjs -v major/minor -n "release notes"
* */

// don't output the commands themselves
$.verbose = true

const versionToBump = getVersionToBump()
const latestReleaseTag = await fetchLatestReleasedVersionNumber()
const latestVersion = removeLeadingV(latestReleaseTag)
const nextVersion = await determineNextVersionNumber(latestVersion)

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
    return (await $`semver bump ${versionToBump} ${previous}`).stdout.trim()
}

async function bumpVersionInFiles() {
    await replaceInFile("README.md", `github.com/seatsio/seatsio-go v${latestVersion}`, `"github.com/seatsio/seatsio-go v${nextVersion}`)
}

async function replaceInFile(filename, latestVersion, nextVersion) {
    return await fs.readFile(filename, 'utf8')
        .then(text => {
            if (text.indexOf(latestVersion) < 0) {
                throw new Error('Not the correct version. Could not find ' + latestVersion + ' in ' + filename)
            }
            return text
        })
        .then(text => text.replace(latestVersion, nextVersion))
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
    return await $`curl -X POST https://pkg.go.dev/fetch/github.com/seatsio/seatsio-go@${newTag}`.catch(error => {
        console.error('pkg.go.dev could not be updated')
        throw error
    })
}
