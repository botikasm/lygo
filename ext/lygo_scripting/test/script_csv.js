(function () {

    var path = $paths.GetAbsolutePath("./data.csv"); // $paths.GetWorkspacePath("./test/data.csv")
    var data = $csv.LoadFromFile(path, true);

    return data;

})();
