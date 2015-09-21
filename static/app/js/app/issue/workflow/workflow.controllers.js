'use strict';

angular.module("TrackerApp.Workflow.controllers", [])
    .controller("WorkflowsController", function($scope, BrowserService, WorkflowService){
        $scope.newItem = function(){
            BrowserService.workflow.add();
        }

        $scope.gridConfig = WorkflowService.gridConfig({ source: WorkflowService.grid });
        $scope.gridParams = {};
    })
    .controller("WorkflowController", function($scope, $routeParams, BrowserService, $timeout, WorkflowService, NotificationTypes,
                                               StatusService, WorkflowStepService, utils, UserService, GroupService){
        $scope.item = { id: $routeParams.id};
        $scope.steps = [];

        var clearMetas = function(){
            $scope.meta = {
                group:{
                    selected:[],
                    unselected:[]
                },
                user:{
                    selected:[],
                    unselected:[]
                }
            }
        };

        var loadStepMembers = function(step){
            clearMetas();
            WorkflowStepService.memberGroups(step.id)
                .then(function(data){
                    $scope.meta.group = data;
                }).then(function(){
                    WorkflowStepService.memberUsers(step.id)
                        .then(function(data){
                            $scope.meta.user = data;
                        });
                });
        };

        var loadStatus = function(){
            StatusService.list($routeParams.id)
                .then(function(data){
                    $scope.statuses = data;
                });
        };

        var loadSteps = function(){
            WorkflowStepService.list($routeParams.id)
                .then(function(data){
                    $scope.steps = data;
                });
        };

        if($routeParams.id){
            WorkflowService.createMeta($routeParams.id)
                .then(function(data){
                    $scope.item=data.item;
                    $scope.steps = data.steps;
                    $scope.statuses = data.statuses;
                });
        }

        GroupService.list()
            .then(function(data){
                $scope.groups = data;
            })
            .then(function(){
                UserService.list()
                    .then(function(data){
                        $scope.users = data;
                    });
            });


        var dialogs;

        $timeout(function(){
            dialogs = {
                statusDialog: $("#statusDialog"),
                stepDialog: $("#stepDialog")
            }

        }, 300);


        $scope.exit = function(){
            BrowserService.workflow.grid();
        };

        $scope.deleteItem = function(){
            if(!utils.confirm()){ return; }
            WorkflowService.remove($scope.item.id).then(function(){ $scope.exit(); });
        };

        $scope.canDelete = function(){
            return $scope.item.id;
        };

        $scope.canSave = function(){
            return true;
        };

        $scope.saveItem = function(){
            var isNewItem = !$scope.item.id;
            WorkflowService.save($scope.item).then(function(data){
                if(isNewItem){
                    BrowserService.workflow.edit(data.id);
                } else {
                    $scope.exit();
                }
            });
        };

        $scope.addStatus = function(){
            $scope.selectedStatus = { name: "New Status" };
            dialogs.statusDialog.modal("show");
        };

        $scope.editStatus = function(status){
            $scope.selectedStatus = {};
            angular.copy(status, $scope.selectedStatus);
            dialogs.statusDialog.modal("show");
        };

        $scope.saveStatus = function(s){
            s.workflow = { id: $scope.item.id };
            StatusService.save(s)
                .then(function(){
                    dialogs.statusDialog.modal("hide");
                    loadStatus();
                    loadSteps();
                });
        };

        $scope.removeStatus = function(s){
            if(!utils.confirm()){ return; }
            StatusService.remove(s.id)
                .then(function(){
                    dialogs.statusDialog.modal("hide");
                    loadStatus();
                });
        };


        $scope.addStep = function(){
            $scope.selectedStep = { name: "New Step" };
            dialogs.stepDialog.modal("show");
        };

        $scope.editStep = function(step){
            $scope.selectedStep = {};
            angular.copy(step, $scope.selectedStep);
            loadStepMembers(step);
            dialogs.stepDialog.modal("show");
        };

        $scope.saveStep = function(s){
            s.workflow = { id: $scope.item.id };
            WorkflowStepService.save(s)
                .then(function(){
                    dialogs.stepDialog.modal("hide");
                    loadSteps();
                });
        };

        $scope.removeStep = function(s){
            if(!utils.confirm()){ return; }
            WorkflowStepService.remove(s.id)
                .then(function(){
                    dialogs.stepDialog.modal("hide");
                    loadSteps();
                });
        };

        //prueba de workflow
        $scope.testStatuses = [];

        $scope.testWorkflow = function(nextStatus){
            $scope.currentStatus= nextStatus;
            nextStatus = nextStatus || { id: ""};
            WorkflowStepService.availableSteps($scope.item.id, nextStatus.id).then(function(data){
                    $scope.testStatuses = data;
                });
        };

        $scope.addMember = function(item){
            var data = {
                workflowStep: { id: $scope.selectedStep.id }
            };

            if(item.meta.documentType === NotificationTypes.objectType.group){
                data.group = item;
            } else if (item.meta.documentType === NotificationTypes.objectType.user){
                data.user = item;
            }
            WorkflowStepService.memberAdd(data)
                .then(function(data){
                    if(item.meta.documentType === NotificationTypes.objectType.group){
                        $scope.meta.group = data;
                    } else if (item.meta.documentType === NotificationTypes.objectType.user){
                        $scope.meta.user = data;
                    }
                });
        };

        $scope.removeMember = function(item){
            var data = {
                workflowStep: { id: $scope.selectedStep.id }
            };

            if(item.meta.documentType === NotificationTypes.objectType.group){
                data.group = item;
            } else if (item.meta.documentType === NotificationTypes.objectType.user){
                data.user = item;
            }
            WorkflowStepService.memberRemove(data)
                .then(function(data){
                    if(item.meta.documentType === NotificationTypes.objectType.group){
                        $scope.meta.group = data;
                    } else if (item.meta.documentType === NotificationTypes.objectType.user){
                        $scope.meta.user = data;
                    }
                });
        };

    });
