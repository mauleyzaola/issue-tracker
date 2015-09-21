'use strict';

angular.module("TrackerApp.Priority.services", [])
    .factory("PriorityService", function($http, BrowserService, PathService, NotificationTypes,
                                         NotificationService, RunApiService, DefaultStyles){
        return {
            load: function(id){
                return $http.get(PathService.priority.load(id))
                    .then(function(response){
                        return response.data;
                    });
            },

            save: function(data){
                var baseFunc = data.id ? $http.put : $http.post;
                return baseFunc(PathService.priority.save, data)
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.priority,
                            operation:data.id ? NotificationTypes.operation.update : NotificationTypes.operation.add,
                            item:response.data
                        });
                        return response.data;
                    });
            },

            remove: function(id){
                return $http.delete(PathService.priority.remove(id))
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.priority,
                            operation:NotificationTypes.operation.delete,
                            item:response.data
                        });
                        return response.data;
                    });
            },

            grid: function(data){
                return $http.get(RunApiService.generateUrl(PathService.priority.grid, data))
                    .then(function(response){
                        return response.data;
                    });
            },


            gridConfig: function(data){
                data = data || {};
                var pars = {
                    customCss: DefaultStyles.css.defaultTableHoverCss,
                    columns: [
                        { name: "Name", field:"name" },
                        { name: "Last Change", field:"lastModified", filter:"timeAgo" }
                    ],
                    rowClick: function(row){
                        BrowserService.priority.edit(row.id);
                    }
                };
                return angular.extend(pars, data);
            },


            list: function(){
                return $http.get(PathService.priority.list)
                    .then(function(response){
                        return response.data;
                    });
            }

        }
    })
