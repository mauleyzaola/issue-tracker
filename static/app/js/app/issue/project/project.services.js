'use strict';

angular.module("TrackerApp.Project.services", [])
    .factory("ProjectService", function($http, $location, PathService, NotificationService,
                                        NotificationTypes, RunApiService, DefaultStyles){
        return {
            createMeta: function(id){
                id = id || "";
                return $http.get(PathService.project.createMeta(id))
                    .then(function(response){
                        return response.data;
                    });
            },

            load: function(id){
                return $http.get(PathService.project.load(id))
                    .then(function(response){
                        return response.data;
                    });
            },

            save: function(data){
                var baseFunc = data.id ? $http.put : $http.post,
                    baseEndPoint = data.id ? PathService.project.update : PathService.project.create;
                return baseFunc(baseEndPoint, data)
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.project,
                            operation:data.id ? NotificationTypes.operation.update : NotificationTypes.operation.add,
                            item:response.data
                        });
                        return response.data;
                    });
            },

            remove: function(id){
                return $http.delete(PathService.project.remove(id))
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.project,
                            operation:NotificationTypes.operation.delete,
                            item:response.data
                        });
                        return response.data;
                    });
            },

            grid: function(data){
                return $http.get(RunApiService.generateUrl(PathService.project.grid, data))
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
                        { name: "Key", field:"pkey" },
                        { name: "Lead", field:"projectLead" },
                        { name: "Starts", field:"begins", filter:"dateFormat" },
                        { name: "Ends", field:"ends", filter:"dateFormat" },
                        { name: "Created", field:"dateCreated", filter:"timeAgo" },
                        { name: "Issues", field:"issueCount" },
                        { name: "% Completed", field:"percentageCompleted", filter:"percent" },
                        { name: "Scheme", field:"permissionScheme" }
                    ],
                    rowClick: function(row){
                        $location.search("");
                        $location.path("/issue/project/" + row.id);
                    }
                };
                return angular.extend(pars, data);
            },

            projectRoles:function(id){
                return $http.get(PathService.project.projectRoles(id))
                    .then(function(response){
                        return response.data;
                    });
            },

            projectRoleMembers:function(id){
                return $http.get(PathService.project.projectRoleMembers(id))
                    .then(function(response){
                        return response.data;
                    });
            },

            projectRoleProjectMembers:function(id){
                return $http.get(PathService.project.projectRoleProjectMembers(id))
                    .then(function(response){
                        return response.data;
                    });
            },

            projectRoleMemberAdd:function(data){
                return $http.post(PathService.project.projectRoleMemberAdd,data)
                    .then(function(response){
                        return response.data;
                    });
            },

            projectRoleMemberRemove:function(data){
                return $http.post(PathService.project.projectRoleMemberRemove,data)
                    .then(function(response){
                        return response.data;
                    });
            }
        }
    })