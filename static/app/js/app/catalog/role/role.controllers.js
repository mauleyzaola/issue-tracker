'use strict';

angular.module("TrackerApp.Role.controllers", [])
    .controller("RolesController", function($scope, BrowserService, RoleService){
        $scope.newItem = function(){
            BrowserService.role.add();
        }

        $scope.gridConfig = RoleService.gridConfig( { source: RoleService.grid });
        $scope.gridParams = {};
    })
    .controller("RoleController", function($scope, $routeParams, BrowserService, RoleService, utils){

        $scope.item = { id: $routeParams.id};

        if($routeParams.id){
            RoleService.load($routeParams.id)
                .then(function(data){ $scope.item=data; });
        }

        $scope.exit = function(){
            BrowserService.role.grid();
        }

        $scope.deleteItem = function(){
            if(!utils.confirm()){ return; }
            RoleService.remove($scope.item.id).then(function(){ $scope.exit(); });
        }

        $scope.canDelete = function(){
            return $scope.item.id;
        }

        $scope.canSave = function(){
            return true;
        }

        $scope.saveItem = function(){
            RoleService.save($scope.item).then(function(data){
                $scope.exit();
            });
        }
    });
