'use strict';

angular.module("TrackerApp.Workflow.services", [])
    .factory("WorkflowService", function($http, BrowserService, PathService, NotificationService,
                                         NotificationTypes, RunApiService, DefaultStyles){
        return {
            load: function(id){
                return $http.get(PathService.workflow.load(id))
                    .then(function(response){
                        return response.data;
                    });
            },

            createMeta: function(id){
                return $http.get(PathService.workflow.createMeta(id))
                    .then(function(response){
                        return response.data;
                    });
            },

            save: function(data){
                var baseFunc = data.id ? $http.put : $http.post;
                return baseFunc(PathService.workflow.save, data)
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.workflow,
                            operation:data.id ? NotificationTypes.operation.update : NotificationTypes.operation.add,
                            item:response.data
                        });
                        return response.data;
                    });
            },

            remove: function(id){
                return $http.delete(PathService.workflow.remove(id))
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.workflow,
                            operation:NotificationTypes.operation.delete,
                            item:response.data
                        });
                        return response.data;
                    });
            },

            grid: function(data){
                return $http.get(RunApiService.generateUrl(PathService.workflow.grid, data))
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
                        BrowserService.workflow.edit(row.id);
                    }
                };
                return angular.extend(pars, data);
            },


            list: function(){
                return $http.get(PathService.workflow.list)
                    .then(function(response){
                        return response.data;
                    });
            }

        }
    })
    .factory("WorkflowStepService", function($http, PathService, NotificationService, NotificationTypes){
        return {
            save: function(data){
                var baseFunc = data.id ? $http.put : $http.post;
                return baseFunc(PathService.workflowStep.save, data)
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.workflowStep,
                            operation:data.id ? NotificationTypes.operation.update : NotificationTypes.operation.add,
                            item:response.data
                        });
                        return response.data;
                    });
            },

            remove: function(id){
                return $http.delete(PathService.workflowStep.remove(id))
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.workflowStep,
                            operation:NotificationTypes.operation.delete,
                            item:response.data
                        });
                        return response.data;
                    });
            },


            list: function(workflow){
                return $http.get(PathService.workflowStep.list(workflow))
                    .then(function(response){
                        return response.data;
                    });
            },

            availableSteps: function(workflow, status){
                return $http.get(PathService.workflowStep.availableSteps({workflow:workflow, status:status}))
                    .then(function(response){
                        return response.data;
                    });
            },

            availableStepsUser: function(workflow, status){
                return $http.get(PathService.workflowStep.availableStepsUser({workflow:workflow, status:status}))
                    .then(function(response){
                        return response.data;
                    });
            },

            members:function(id){
                return $http.get(PathService.workflowStep.members(id))
                    .then(function(response){
                        return response.data;
                    });
            },

            memberAdd:function(data){
                return $http.post(PathService.workflowStep.memberAdd,data)
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.workflowStep,
                            operation:NotificationTypes.operation.add,
                            item:response.data
                        });
                        return response.data;
                    });
            },

            memberRemove:function(data){
                return $http.post(PathService.workflowStep.memberRemove,data)
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.workflowStep,
                            operation:NotificationTypes.operation.delete,
                            item:response.data
                        });
                        return response.data;
                    });
            },

            memberGroups:function(id){
                return $http.get(PathService.workflowStep.memberGroups(id))
                    .then(function(response){
                        return response.data;
                    });
            },

            memberUsers:function(id){
                return $http.get(PathService.workflowStep.memberUsers(id))
                    .then(function(response){
                        return response.data;
                    });
            }

        }
    })
    .factory("StatusService", function($http, PathService, NotificationService, NotificationTypes){
        return {
            load: function(id){
                return $http.get(PathService.status.load(id))
                    .then(function(response){
                        return response.data;
                    });
            },

            save: function(data){
                var baseFunc = data.id ? $http.put : $http.post;
                return baseFunc(PathService.status.save, data)
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.status,
                            operation:data.id ? NotificationTypes.operation.update : NotificationTypes.operation.add,
                            item:response.data
                        });
                        return response.data;
                    });
            },

            remove: function(id){
                return $http.delete(PathService.status.remove(id))
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.status,
                            operation:NotificationTypes.operation.delete,
                            item:response.data
                        });
                        return response.data;
                    });
            },


            list: function(workflow){
                return $http.get(PathService.status.list(workflow))
                    .then(function(response){
                        return response.data;
                    });
            }

        }
    })