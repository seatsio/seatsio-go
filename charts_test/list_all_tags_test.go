package charts

/* TODO
func TestListAllTags(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	chartKey1 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client.Charts.AddTag(chartKey1, "tag1")
	client.Charts.AddTag(chartKey1, "tag2")

	chartKey2 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client.Charts.AddTag(chartKey2, "tag3")

	list, err := client.Charts.ListAllTags()
	require.NoError(t, err)

	require.Equal(t, 3, len(list))
	require.Contains(t, list, "tag1")
	require.Contains(t, list, "tag2")
	require.Contains(t, list, "tag3")
}
*/
