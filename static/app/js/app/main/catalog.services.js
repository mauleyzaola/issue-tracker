'use strict';

angular.module("TrackerApp.Catalog.services", [])
    .factory("FileOperations", function($http, Upload, PathService, Notifier){
        return {

            fileSelect: function($files, callback) {
                var maxFileSize = 1024 * 1024 * 10;

                for (var i = 0; i < $files.length; i++) {
                    var $file = $files[i];

                    if($file.size > maxFileSize){
                        Notifier.error({
                            title: "File too large",
                            text: $file.name + " is larger than maximum allowed."
                        });
                        return;
                    }


                    Upload.upload({
                        url: PathService.file.upload,
                        method: "post",
                        file: $file
                    }).success(function(data) {
                        if(typeof(callback) === "function"){
                            callback(data);
                        }
                    });
                }
            }
        }
    })