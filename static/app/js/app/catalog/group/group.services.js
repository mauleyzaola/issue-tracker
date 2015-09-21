'use strict';

angular.module("TrackerApp.Group.services", [])
    .factory("GroupService", function($http, BrowserService, PathService, NotificationTypes,
                                         NotificationService, RunApiService, DefaultStyles){
        return {
            load: function(id){
                return $http.get(PathService.group.load(id))
                    .then(function(response){
                        return response.data;
                    });
            },

            save: function(data){
                var baseFunc = data.id ? $http.put : $http.post;
                return baseFunc(PathService.group.save, data)
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.group,
                            operation:data.id ? NotificationTypes.operation.update : NotificationTypes.operation.add,
                            item:response.data
                        });
                        return response.data;
                    });
            },

            remove: function(id){
                return $http.delete(PathService.group.remove(id))
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.group,
                            operation:NotificationTypes.operation.delete,
                            item:response.data
                        });
                        return response.data;
                    });
            },

            grid: function(data){
                return $http.get(RunApiService.generateUrl(PathService.group.grid, data))
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
                        { name: "Created", field:"dateCreated", filter:"timeAgo" },
                        { name: "Last Changed", field:"lastModified", filter:"timeAgo" }
                    ],
                    rowClick: function(row){
                        BrowserService.group.edit(row.id);
                    }
                };
                return angular.extend(pars, data);
            },

            list: function(){
                return $http.get(PathService.group.list)
                    .then(function(response){
                        return response.data;
                    });
            },

            groups:function(id){
                return $http.get(PathService.group.groups(id))
                    .then(function(response){
                        return response.data;
                    });
            },

            users:function(id){
            return $http.get(PathService.group.users(id))
                .then(function(response){
                    return response.data;
                });
            },

            addGroupUser:function(data){
                return $http.post(PathService.group.addGroupUser, data)
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.group,
                            operation:NotificationTypes.operation.add
                        });
                        return response.data;
                    })
            },

            removeGroupUser:function(data){
                return $http.post(PathService.group.removeGroupUser, data)
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.group,
                            operation:NotificationTypes.operation.delete
                        });
                        return response.data;
                    })
            }

        }
    })
