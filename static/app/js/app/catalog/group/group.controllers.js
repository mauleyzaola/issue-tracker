'use strict';

angular.module("TrackerApp.Group.controllers", [])
    .controller("GroupsController", function($scope, BrowserService, GroupService){
        $scope.newItem = function(){
            BrowserService.group.add();
        }

        $scope.gridConfig = GroupService.gridConfig( { source: GroupService.grid });
        $scope.gridParams = {};
    })
    .controller("GroupController", function($scope, $routeParams, BrowserService, GroupService, utils){

        $scope.item = { id: $routeParams.id};
        $scope.users={
            selected:[],
            unselected:[]
        };

        if($routeParams.id){
            GroupService.load($routeParams.id)
                .then(function(data){ $scope.item=data; })
                .then(function(){
                    GroupService.users($routeParams.id)
                        .then(function(data){
                            $scope.users = data;
                        });
                });
        }



        $scope.exit = function(){
            BrowserService.group.grid();
        }

        $scope.deleteItem = function(){
            if(!utils.confirm()){ return; }
            GroupService.remove($scope.item.id).then(function(){ $scope.exit(); });
        }

        $scope.canDelete = function(){
            return $scope.item.id;
        }

        $scope.canSave = function(){
            return true;
        }

        $scope.saveItem = function(){
            var isNewItem = !$scope.item.id;
            GroupService.save($scope.item).then(function(data){
                if(!isNewItem) { $scope.exit(); }
                else {
                    BrowserService.group.edit(data.id);
                }
            });
        }

        $scope.addUser = function(index){
            var selected = $scope.users.unselected[index];
            var data = {
                user: { id: selected.id },
                group: { id: $scope.item.id }
            }
            GroupService.addGroupUser(data)
                .then(function(){
                    $scope.users.unselected.splice(index, 1);
                    $scope.users.selected.push(selected);
                })
        }

        $scope.removeUser = function(index){
            var selected = $scope.users.selected[index];
            var data = {
                user: { id: selected.id },
                group: { id: $scope.item.id }
            }
            GroupService.removeGroupUser(data)
                .then(function(){
                    $scope.users.selected.splice(index, 1);
                    $scope.users.unselected.push(selected);
                })
        }

    });
