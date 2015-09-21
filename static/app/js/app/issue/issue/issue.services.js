'use strict';

angular.module("TrackerApp.Issue.services", [])
    .factory("IssueService", function($http, BrowserService, PathService, PathResolver, NotificationService,
                                      NotificationTypes, RunApiService, DefaultStyles){
        return {
            createMeta: function(data){
                return $http.get(PathService.issue.createMeta(data))
                    .then(function(response){
                        return response.data;
                    });
            },

            load: function(data){
                return $http.get(PathService.issue.load(data))
                    .then(function(response){
                        return response.data;
                    });
            },

            move:function(data){
                return $http.post(PathService.issue.move(data.id), data)
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.issue,
                            operation:NotificationTypes.operation.update,
                            item:response.data
                        });
                        return response.data;
                    });
            },

            save: function(data){
                var baseFunc = data.id ? $http.put : $http.post;
                return baseFunc(PathService.issue.save(data), data)
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.issue,
                            operation:data.id ? NotificationTypes.operation.update : NotificationTypes.operation.add,
                            item:response.data
                        });
                        return response.data;
                    });
            },

            remove: function(id){
                return $http.delete(PathService.issue.remove(id))
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.issue,
                            operation:NotificationTypes.operation.delete,
                            item:response.data
                        });
                        return response.data;
                    });
            },

            assignToMe:function(data){
                return $http.post(PathService.issue.assignToMe(data))
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.issue,
                            operation:NotificationTypes.operation.update,
                            item:response.data
                        });
                        return response.data;
                    });
            },

            reporterIsMe:function(data){
                return $http.post(PathService.issue.reporterIsMe(data))
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.issue,
                            operation:NotificationTypes.operation.update,
                            item:response.data
                        });
                        return response.data;
                    });
            },

            grid: function(data){
                return $http.get(RunApiService.generateUrl(PathService.issue.grid, data))
                    .then(function(response){
                        return response.data;
                    });
            },


            gridConfig: function(data){
                data = data || {};
                var pars = {
                    customCss: DefaultStyles.css.defaultTableHoverCss,
                    columns: [
                        { name: "Key", field:"pkey" },
                        { name: "Project", field:"project" },
                        { name: "Parent", field:"parent" },
                        { name: "Name", field:"name" },
                        { name: "Assignee", field:"assignee" },
                        { name: "Reporter", field:"reporter" },
                        { name: "Priority", field:"priority" },
                        { name: "Status", field:"status" },
                        { name: "Due", field:"dueDate", filter:"timeAgo" }
                    ],
                    rowClick: function(row){
                        BrowserService.issue.edit(row.pkey);
                    }
                };
                return angular.extend(pars, data);
            },

            currentUserSubscribed:function(id){
                return $http.get(PathService.issue.currentUserSubscribed(id))
                    .then(function(response){
                        return response.data;
                    });
            },

            changeStatus: function(pkey, status){
                return $http.post(PathService.issue.changeStatus({pkey:pkey, status:status}))
                .then(function(response){
                    NotificationService.notify({
                        objectType: NotificationTypes.objectType.issue,
                        operation:NotificationTypes.operation.update,
                        item:response.data
                    });
                    return response.data;
                });
            },

            commentAdd:function(data){
                return $http.post(PathService.issue.commentAdd(data),data)
                    .then(function(response){
                        NotificationService.notify({
                            objectType:NotificationTypes.objectType.issueComment,
                            operation:NotificationTypes.operation.add,
                            item: response.data
                        });
                        return response.data;
                    });
            },

            commentUpdate:function(data){
                return $http.put(PathService.issue.commentUpdate(data),data)
                    .then(function(response){
                        NotificationService.notify({
                            objectType:NotificationTypes.objectType.issueComment,
                            operation:NotificationTypes.operation.update,
                            item: response.data
                        });
                        return response.data;
                    });
            },

            commentRemove:function(data){
                return $http.delete(PathService.issue.commentRemove(data))
                    .then(function(response){
                        NotificationService.notify({
                            objectType:NotificationTypes.objectType.issueComment,
                            operation:NotificationTypes.operation.delete,
                            item: response.data
                        });
                        return response.data;
                    });
            },

            commentList:function(data){
                return $http.get(PathService.issue.commentList(data))
                    .then(function(response){
                        return response.data;
                    });
            },

            subscriptionToggle:function(data){
                return $http.post(PathService.issue.subscriptionToggle(data))
                    .then(function(response){
                        return response.data;
                    });
            },

            subscriptionToggleAny:function(data){
                return $http.post(PathService.issue.subscriptionToggleAny,data)
                    .then(function(response){
                        return response.data;
                    });
            },

            subscribedSelected:function(id){
                return $http.get(PathService.issue.subscribedSelected(id))
                    .then(function(response){
                        return response.data;
                    });
            },

            mySubscriptionsGrid:function(data){
                return $http.get(RunApiService.generateUrl(PathService.issue.mySubscriptions,data))
                    .then(function(response){
                        return response.data;
                    });
            },

            attachments:function(data){
                return $http.get(PathService.issue.attachments(data))
                    .then(function (response) {
                        return response.data;
                    });
            },

            attachmentLoad:function(id){
                return $http.get(PathService.issue.attachmentLoad(id))
                    .then(function(response){
                        return response.data;
                    });
            },

            attachmentAdd:function(data){
                return $http.post(PathService.issue.attachmentAdd(data), data)
                    .then(function(response){
                        return response.data;
                    });
            },

            attachmentRemove:function(data){
                return $http.delete(PathService.issue.attachmentRemove(data))
                    .then(function(response){
                        return response.data;
                    });
            },

            getChildren:function(id){
                return $http.get(PathService.issue.getChildren(id))
                    .then(function(response){
                        return response.data;
                    });
            },

            getDocument:function(id){
                return $http.get(PathService.issue.getDocument(id))
                    .then(function(response){
                        return response.data;
                    });
            }
        }
    })