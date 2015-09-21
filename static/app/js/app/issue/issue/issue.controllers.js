'use strict';

angular.module("TrackerApp.Issue.controllers", [])
    .controller('IssueAttachmentForwarderController', function($location, $routeParams, IssueService, BrowserService){
        var id = $routeParams.id;
        if(!id){
            $location.$$search = {};
            $location.path('/');
        } else {
            IssueService.attachmentLoad(id)
                .then(function(data){
                    BrowserService.issue.edit(data.issue.pkey);
                });
        }
    })
    .controller("IssuesController", function($scope, $location, BrowserService, IssueService){
        $scope.newItem = function(){
            BrowserService.issue.add();
        }

        $scope.gridConfig = IssueService.gridConfig({
            source: IssueService.grid,
            columns: [
                { name: "Key", field:"pkey" },
                { name: "Project", field:"project" },
                { name: "Belongs to", field:"parent" },
                { name: "Name", field:"name" },
                { name: "Assignee", field:"assignee" },
                { name: "Reporter", field:"reporter" },
                { name: "Priority", field:"priority" },
                { name: "Status", field:"status" },
                { name: "Due Date", field:"dueDate", filter:"timeAgo" }
            ]
        });
        $scope.gridParams = $location.$$search;
    })
    .controller("IssueController", function($scope, $location, $routeParams, BrowserService, BrowserUrlService, NotificationTypes,
                                            PathResolver, $timeout, FileOperations, IssueService, utils, PathService, QueryStringNames,
                                            Notifier, WorkflowStepService,
                                            ProjectService, PermissionSchemeService){

        $scope.item = { id: $routeParams.id};
        $scope.attachments = [];
        $scope.subtasks = [];
        $scope.isEditable = true;
        var permissions = {};

        var dialogs = {
            commentDialog:$("#commentDialog"),
            projectSelectDialog:$("#projectSelectDialog")
        };

        var setEditableIssue = function(issue){
            var value = true;

            if(issue.id){
                if(issue.resolvedDate){
                    value = false;
                } else if ($scope.item.cancelledDate){
                    value = false;
                } else if (!permissions.EDIT_ISSUE){
                    value = false;
                }
            }
            $scope.isEditable = value;
            $scope.permissions = permissions;
        }


        var loadSubscriptions = function(){
            IssueService.subscribedSelected($scope.item.id)
                .then(function(data){
                    $scope.subscribers = data;
                });
        };



        var loadIssue = function(){
            var data = {
                pkey:$routeParams.pkey,
                parent:$location.$$search.parent
            };

            IssueService.createMeta(data)
                .then(function(data){
                    $scope.item = data.item;
                    $scope.comments = data.comments;
                    $scope.priorities = data.priorities;
                    $scope.workflows = data.workflows;
                    $scope.steps = data.steps;
                    $scope.reporters = angular.copy(data.users, []);
                    $scope.isSubscribed = data.isSubscribed;
                    data.users.unshift({name:"None"});
                    $scope.assignees = angular.copy(data.users, []);

                    if($routeParams.pkey){
                        $scope.item.parent = data.parent;

                        IssueService.getChildren($scope.item.id)
                            .then(function(data){
                                $scope.subtasks=data;
                            }).then(function(){
                                if($scope.attachments.length == 0){
                                    IssueService.attachments($scope.item)
                                        .then(function(data){
                                            $scope.attachments=data;
                                        });
                                }
                            }).then(function(){
                                PermissionSchemeService.permissionAvailableUser({
                                    user:{id:$scope.currentSession.user.id},
                                    issue:{id:$scope.item.id, project:$scope.item.project}
                                }).then(function(data){
                                    permissions = {};
                                    angular.forEach(data, function(p){
                                        permissions[p.name] = true;
                                    });
                                    setEditableIssue($scope.item);
                                });
                            }).then(function(){
                                loadSubscriptions();
                            })
                    } else {
                        var idProject = $location.$$search[QueryStringNames.project];
                        if(idProject){
                            ProjectService.load(idProject)
                                .then(function(data){
                                    $scope.item.project = data;
                                });
                        }
                        var idParent = $location.$$search[QueryStringNames.parent];
                        if (idParent){
                            IssueService.load({id:idParent})
                                .then(function(data){
                                    $scope.item.parent = data;
                                    $scope.item.project = data.project;
                                });
                        }

                        PermissionSchemeService.names()
                            .then(function(data){
                                permissions = {};
                                angular.forEach(data, function(p){
                                    permissions[p.name] = true;
                                });
                                setEditableIssue($scope.item);
                            });
                    }
                });

            if(!$scope.gridProject){
                $scope.gridProject = ProjectService.gridConfig({
                    source:ProjectService.grid,
                    rowClick:function(r){
                        var oldId = '';
                        if($scope.item.project){ oldId = $scope.item.project.id; }
                        if(oldId === r.id){
                            dialogs.projectSelectDialog.modal('hide');
                            return;
                        }
                        if($scope.item.id){
                            if(!utils.confirm('Sure to move this issue to project :' + r.name)){
                                dialogs.projectSelectDialog.modal('hide');
                                return;
                            }
                            IssueService.move({
                                project:{id: r.id},
                                id:$scope.item.id
                            }).then(function(data){
                                dialogs.projectSelectDialog.modal('hide');
                                $timeout(function(){
                                    BrowserService.issue.edit(data.pkey);
                                },100);
                            });
                        } else {
                            dialogs.projectSelectDialog.modal('hide');
                            ProjectService.load(r.id)
                                .then(function(data){
                                    $scope.item.project = data;
                                });
                        }
                    }
                }) ;

                $scope.gridProjectParams = { };
            }
        }

        loadIssue();

        $scope.exit = function(){
            BrowserService.issue.grid();
        }

        $scope.deleteItem = function(){
            if(!utils.confirm()){ return; }
            IssueService.remove($scope.item.id).then(function(){ $scope.exit(); });
        }

        $scope.canDelete = function(){
            return $scope.item.id && $scope.subtasks.length == 0 && (!$scope.document || $scope.document.meta == null)
                && permissions.DELETE_ISSUE;
        }

        $scope.canSave = function(){
            return $scope.isEditable;
        }

        $scope.saveItem = function(){
            IssueService.save($scope.item).then(function(data){
                $scope.exit();
            });
        }

        $scope.changeStep = function(step){
            IssueService.changeStatus($scope.item.pkey, step.nextStatus.id)
                .then(function(){
                    loadIssue();
                });
        }

        $scope.addComment = function(){
            var item={
                dateCreated:new Date(),
                user:$scope.currentSession.user
            }
            $scope.selectedComment = item;
        }

        $scope.editComment = function(c){
            var item = {};
            angular.copy(c,item);
            $scope.selectedComment = item;
        }

        $scope.saveComment = function(c){
            var index = utils.findIndexById(c.id, $scope.comments);
            c.issue = $scope.item;
            var baseFunc = index == -1 ? IssueService.commentAdd : IssueService.commentUpdate;
            baseFunc(c).then(function(data){
                if(index == -1){
                    $scope.comments.push(data);
                } else {
                    $scope.comments[index]=data;
                }
                dialogs.commentDialog.modal("hide");
            });
        }

        $scope.removeComment = function(c){
            if(!utils.confirm("Sure to remove this comment?")){ return; }
            var index=utils.findIndexById(c.id, $scope.comments);
            c.issue=$scope.item;
            IssueService.commentRemove(c)
                .then(function(){
                    dialogs.commentDialog.modal("hide");
                    $scope.comments.splice(index,1);
                })
        }

        $scope.toggleSubscription = function(){
            IssueService.subscriptionToggle($scope.item)
                .then(function(data){
                    $scope.isSubscribed = data.selected;
                });
        }

        $scope.onFileSelect = FileOperations.fileSelect;

        $scope.addAttachment = function(data){
            IssueService.attachmentAdd({
                issue:$scope.item,
                fileItem:data
            }).then(function(data){
                $scope.attachments.push(data);
            });
        }


        $scope.deleteAttachment = function(a, index){
            if(!utils.confirm("Sure to remove the attachment?")) {return;}
            IssueService.attachmentRemove(a)
                .then(function(data){
                    $scope.attachments.splice(index,1);
                });
        }

        $scope.browseIssue = function(i){
            BrowserService.issue.edit(i.pkey);
        }

        $scope.addSubtask = function(){
            BrowserService.issue.add({parent:$scope.item.id});
        }

        $scope.resolveIssueUrl=function(issue){
            if(!issue){ return; }
            $location.path(BrowserUrlService.issue.edit(issue.pkey));
        }

        $scope.browsePo=function(po){
            if (!po){return;}
            return BrowserUrlService.purchaseOrder.edit(po.id);
        }

        $scope.browseStockMovement = function(stock){
            if(!stock){return;}
            return BrowserUrlService.stockMovement.edit(stock.stockMovementHeader.id);
        }


        $scope.assignToMe = function(){
            IssueService.assignToMe({id:$scope.item.id})
                .then(function(data){
                    $scope.item.assignee = data.assignee;
                });
        }

        $scope.reporterIsMe= function(){
            IssueService.reporterIsMe({id:$scope.item.id})
                .then(function(data){
                    $scope.item.reporter = data.reporter;
                });
        }

        $scope.toggleSubscriptionAny = function(user){
            IssueService.subscriptionToggleAny({
                issue:{ id:$scope.item.id },
                user: { id:user.id }
            }).then(function(data){
                $scope.subscribers = data;
            });
        }
    });
