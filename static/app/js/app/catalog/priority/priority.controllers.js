'use strict';

angular.module("TrackerApp.Priority.controllers", [])
    .controller("PrioritiesController", function($scope, BrowserService, PriorityService){
        $scope.newItem = function(){
            BrowserService.priority.add();
        }

        $scope.gridConfig = PriorityService.gridConfig({ source: PriorityService.grid });
        $scope.gridParams = {};
    })
    .controller("PriorityController", function($scope, $routeParams, BrowserService, PriorityService, utils){
        $scope.item = { id: $routeParams.id};

        if($routeParams.id){
            PriorityService.load($routeParams.id)
                .then(function(data){ $scope.item=data;  });
        }

        $scope.exit = function(){
            BrowserService.priority.grid();
        }

        $scope.deleteItem = function(){
            if(!utils.confirm()){ return; }
            PriorityService.remove($scope.item.id).then(function(){ $scope.exit(); });
        }

        $scope.canDelete = function(){
            return $scope.item.id;
        }

        $scope.canSave = function(){
            return true;
        }

        $scope.saveItem = function(){
            PriorityService.save($scope.item).then(function(){ $scope.exit(); });
        }

    });
