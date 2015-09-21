'use strict';

angular.module("TrackerApp.User.services", [])
    .factory("UserService", function($http, BrowserService, PathService, NotificationTypes,
                                        NotificationService, RunApiService, DefaultStyles){
        return {
            load: function(id){
                return $http.get(PathService.user.load(id))
                    .then(function(response){
                        return response.data;
                    });
            },

            save: function(data){
                var baseFunc = data.id ? $http.put : $http.post;
                return baseFunc(PathService.user.save, data)
                        .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.user,
                            operation:data.id ? NotificationTypes.operation.update : NotificationTypes.operation.add,
                            item:response.data
                        });
                        return response.data;
                    });
            },

            remove: function(id){
                return $http.delete(PathService.user.remove(id))
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.user,
                            operation:NotificationTypes.operation.delete,
                            item:response.data
                        });
                        return response.data;
                    });
            },

            grid: function(data){
                return $http.get(RunApiService.generateUrl(PathService.user.grid, data))
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
                        { name: "Last Name", field:"lastName" },
                        { name: "Email", field:"email" },
                        { name: "Last Change", field:"lastModified", filter:"timeAgo" },
                        { name: "Last Login", field:"lastLogin", filter: "timeAgo" }
                    ],
                    rowClick: function(row){
                        BrowserService.user.edit(row.id);
                    }
                };
                return angular.extend(pars, data);
            },

            list: function(){
                return $http.get(PathService.user.list)
                    .then(function(response){
                        return response.data;
                    });
            }

        }
    })
