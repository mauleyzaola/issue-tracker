'use strict';

angular.module("TrackerApp.Role.services", [])
    .factory("RoleService", function($http, BrowserService, PathService, NotificationTypes,
                                        NotificationService, RunApiService, DefaultStyles){
        return {
            load: function(id){
                return $http.get(PathService.role.load(id))
                    .then(function(response){
                        return response.data;
                    });
            },

            save: function(data){
                var baseFunc = data.id ? $http.put : $http.post;
                return baseFunc(PathService.role.save, data)
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.role,
                            operation:data.id ? NotificationTypes.operation.update : NotificationTypes.operation.add,
                            item:response.data
                        });
                        return response.data;
                    });
            },

            remove: function(id){
                return $http.delete(PathService.role.remove(id))
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.role,
                            operation:NotificationTypes.operation.delete,
                            item:response.data
                        });
                        return response.data;
                    });
            },

            grid: function(data){
                return $http.get(RunApiService.generateUrl(PathService.role.grid, data))
                    .then(function(response){
                        return response.data;
                    });
            },

            gridConfig: function(data){
                data = data || {};
                var pars = {
                    customCss: DefaultStyles.css.defaultTableHoverCss,
                    columns: [
                        { name: "Name", field:"name" }
                    ],
                    rowClick: function(row){
                        BrowserService.role.edit(row.id);
                    }
                };
                return angular.extend(pars, data);
            },

            list: function(){
                return $http.get(PathService.role.list)
                    .then(function(response){
                        return response.data;
                    });
            }        }
    })
