'use strict';

angular.module("TrackerApp.User.controllers", [])
    .controller("UsersController", function($scope, BrowserService, UserService){
        $scope.newItem = function(){
            BrowserService.user.add();
        }

        $scope.gridConfig = UserService.gridConfig( { source: UserService.grid });
        $scope.gridParams = {};
    })
    .controller("UserController", function($scope, $routeParams, BrowserService, UserService, AccountService, GroupService, utils, IssueService){

        $scope.item = { id: $routeParams.id};
        $scope.units={
            selected:[],
            unselected:[]
        };

        $scope.groups = {
            selected:[],
            unselected:[]
        };

        if($routeParams.id){
            UserService.load($routeParams.id)
                .then(function(data){ $scope.item=data; })
                .then(function(){
                    GroupService.groups($scope.item.id)
                        .then(function(data){
                            $scope.groups = data;
                        });
                });

            $scope.gridAssignee = IssueService.gridConfig({
                source:IssueService.grid,
                columns: [
                    { name: "Key", field:"pkey" },
                    { name: "Project", field:"project" },
                    { name: "Name", field:"name" },
                    { name: "Reporter", field:"reporter" },
                    { name: "Priority", field:"priority" },
                    { name: "Status", field:"status" },
                    { name: "Due Date", field:"dueDate", filter:"timeAgo" }
                ]
            });
            $scope.gridAssigneeParams = {
                assignee:$scope.item.id,
                resolved:false
            };

            $scope.gridReporter = IssueService.gridConfig({
                source:IssueService.grid,
                columns: [
                    { name: "Key", field:"pkey" },
                    { name: "Project", field:"project" },
                    { name: "Name", field:"name" },
                    { name: "Assignee", field:"assignee" },
                    { name: "Priority", field:"priority" },
                    { name: "Status", field:"status" },
                    { name: "Due Date", field:"dueDate", filter:"timeAgo" }
                ]
            });
            $scope.gridReporterParams = {
                reporter:$scope.item.id,
                resolved:false
            };
        }


        $scope.changePassword = function(){
            AccountService.changePassword($scope.item.id, $scope.item.password1)
                .then(function(){
                    $("#changePasswordDialog").modal("hide");
                    $scope.item.password1=null;
                    $scope.item.password2=null;
                });
        }

        $scope.enableChangePassword = function(){
            return $scope.item.password1 && ($scope.item.password1 === $scope.item.password2);
        }

        $scope.exit = function(){
            BrowserService.user.grid();
        }

        $scope.deleteItem = function(){
            if(!utils.confirm()){ return; }
            UserService.remove($scope.item.id).then(function(){ $scope.exit(); });
        }

        $scope.canDelete = function(){
            return $scope.item.id;
        }

        $scope.canSave = function(){
            return true;
        }

        $scope.saveItem = function(){
            var isNewItem = !$scope.item.id;
            UserService.save($scope.item).then(function(data){
                if(!isNewItem) { $scope.exit(); }
                else {
                    BrowserService.user.edit(data.id);
                }
            });
        }

        $scope.addGroup = function(index){
            var selected = $scope.groups.unselected[index];
            var data = {
                group: { id: selected.id },
                user: { id: $scope.item.id }
            }
            GroupService.addGroupUser(data)
                .then(function(){
                    $scope.groups.unselected.splice(index, 1);
                    $scope.groups.selected.push(selected);
                })
        }

        $scope.removeGroup = function(index){
            var selected = $scope.groups.selected[index];
            var data = {
                group: { id: selected.id },
                user: { id: $scope.item.id }
            }
            GroupService.removeGroupUser(data)
                .then(function(){
                    $scope.groups.selected.splice(index, 1);
                    $scope.groups.unselected.push(selected);
                })
        }

    });
