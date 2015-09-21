'use strict';

angular.module('TrackerApp.FileItem.services', [])
    .factory('FileItemService', function($http, $log, $window, BrowserService, PathService,
                                         RunApiService, DefaultStyles, SessionManagement){
        return {
            directoryGrid: function(data){
                return $http.get(RunApiService.generateUrl(PathService.file.directoryGrid, data))
                    .then(function(response){
                        return response.data;
                    });
            },

            directoryGridConfig: function(data){
                data = data || {};
                var pars =  {
                    customCss: DefaultStyles.css.defaultTableHoverCss,
                    columns: [
                        { name: "Fecha", field:"dateCreated" },
                        { name: "Num. Archivos", field:"itemCount", template: '<span class="pull-right">@</span> ' },
                        { name: "Espacio", field:"bytes", filter:"sizeGeneric", template: '<span class="pull-right">@</span> ' }
                    ],
                    rowClick: function(){
                        $log.log('not implemented on the consumer');
                    }
                };
                return angular.extend(pars, data);
            },

            fileGrid: function(data){
                return $http.get(RunApiService.generateUrl(PathService.file.fileGrid, data))
                    .then(function(response){
                        return response.data;
                    });
            },

            fileGridConfig: function(data){
                data = data || {};
                var pars =  {
                    customCss: DefaultStyles.css.defaultTableHoverCss,
                    columns: [
                        { name: "Fecha", field:"dateCreated", filter:'dateTimeFormat' },
                        { name: "Nombre", field:"name" },
                        { name: "Espacio", field:"bytes", filter:"sizeGeneric", template: '<span class="pull-right">@</span> ' }
                    ],
                    rowClick: function(r){
                        var data = {
                                token:SessionManagement.currentSession().id
                            },
                            url = PathService.file.download + '/' + r.id;
                        $window.open(RunApiService.generateUrl(url, data));
                    }
                };
                return angular.extend(pars, data);
            }

        }
    })
