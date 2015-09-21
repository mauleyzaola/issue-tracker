'use strict';

angular.module("TrackerApp.PermissionScheme.services", [])
    .factory("PermissionNameService", function($http, PathService){

    })
    .factory("PermissionSchemeService", function($http, BrowserService, PathService, NotificationTypes,
                                                 NotificationService, RunApiService, DefaultStyles){
        return {

            names:function(){
                return $http.get(PathService.permissionScheme.names)
                    .then(function(response){
                        return response.data;
                    });
            },

            load: function(id){
                return $http.get(PathService.permissionScheme.load(id))
                    .then(function(response){
                        return response.data;
                    });
            },

            save: function(data){
                var baseFunc = data.id ? $http.put : $http.post;
                return baseFunc(PathService.permissionScheme.save, data)
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.permissionScheme,
                            operation:data.id ? NotificationTypes.operation.update : NotificationTypes.operation.add,
                            item:response.data
                        });
                        return response.data;
                    });
            },

            remove: function(id){
                return $http.delete(PathService.permissionScheme.remove(id))
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.permissionScheme,
                            operation:NotificationTypes.operation.delete,
                            item:response.data
                        });
                        return response.data;
                    });
            },

            clear: function(id){
                return $http.post(PathService.permissionScheme.clear(id))
                    .then(function(response){
                        return response.data;
                    });
            },

            grid: function(data){
                return $http.get(RunApiService.generateUrl(PathService.permissionScheme.grid, data))
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
                        BrowserService.permissionScheme.edit(row.id);
                    }
                };
                return angular.extend(pars, data);
            },

            list:function(){
                return $http.get(PathService.permissionScheme.list)
                    .then(function(response){
                        return response.data;
                    });
            },

            itemAdd:function(data){
                return $http.post(PathService.permissionSchemeItem.add,data)
                    .then(function(response){
                        return response.data;
                    });
            },

            itemRemove:function(data){
                return $http.post(PathService.permissionSchemeItem.remove, data)
                    .then(function(response){
                        return response.data;
                    });
            },

            items:function(id){
                return $http.get(PathService.permissionSchemeItem.list(id))
                    .then(function(response){
                        return response.data;
                    });
            },

            projects:function(id){
                return $http.get(PathService.permissionScheme.projects(id))
                    .then(function(response){
                        return response.data;
                    });
            },

            permissionAvailableUser:function(data){
                return $http.post(PathService.permissionSchemeItem.permissionAvailableUser,data)
                    .then(function(response){
                        return response.data;
                    });
            }
        }
    })
