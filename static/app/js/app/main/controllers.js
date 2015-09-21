'use strict';

var menus = function(){
    var cat = {
        iconCss:'fa-book',
        name:'Catalog',
        items:[
            { name:'Priorities', target:'/catalog/priorities'},
            { name:'Workflows', target:'/catalog/workflows'}
        ]
    };

    var sec = {
        iconCss:'fa-lock',
        name:'Security',
        items:[
            { name:'Groups', target:'/catalog/groups'},
            { name:'Permission Schemes', target:'/catalog/permissionschemes'},
            { name:'Roles', target:'/catalog/roles'},
            { name:'Users', target:'/catalog/users'}
        ]
    };

    var issue = {
        iconCss:'fa-bar-chart',
        name:'Issues',
        items:[
            { name:'Todo', target:'/issue/issues/list' },
            { name:'Resolved', target:'/issue/issues/list' },
            { name:'Projects', target:'/issue/projects' }
        ]
    };

    return [
        cat,
        issue,
        sec
    ];
};


angular.module('TrackerApp.controllers', [])
.controller('AppCtrl', function ($rootScope, $scope, $location, $window, SessionManagement, AccountService, ResolveUrlService) {

        $scope.language = { language: "en" };

        $rootScope.currentSession = SessionManagement.currentSession();

        $scope.defaultTableCss = 'table table-bordered table-striped table-hover';

        $scope.resolveUrl = function(item){
            return ResolveUrlService.resolveUrl(item);
        };

        $scope.isLoggedIn = function(){
            return $rootScope.currentSession && $rootScope.currentSession.id;
        };

        $scope.isSystemAdministrator = function(){
            return $scope.currentSession && $scope.currentSession.user && $scope.currentSession.user.isSystemAdministrator;
        };

        $scope.$on("user:login", function(data, args){
            SessionManagement.currentSession(args);
            $rootScope.currentSession = args;

            if($location.$$search.returnurl){
                var params = $location.$$search.returnurl.split("?");
                $location.path(params[0]);
                if(params[1]) $location.search(params[1]);
            } else {
                $location.path("/");
                $location.search("");
            }

        });

        $scope.$on("user:logout", function(){
            $rootScope.currentSession = null;
            SessionManagement.clearCurrentSession();
            $location.path("/");
            $location.$$search = {};
        });

        $scope.$on('user:update',function(data,args){
            if(!$scope.isLoggedIn()){ return; }
            if($scope.currentSession.user.id === args.id){
                $scope.currentSession.user = args;
                SessionManagement.currentSession($scope.currentSession);
            }
        });

        $scope.menuItems = menus();
    })
    .controller("Index.Controller", function($scope, $timeout, DashboardService, NotificationTypes,
                                             BrowserService, BrowserUrlService, RunApiService){

        $scope.meta = {
            my:{
                assignee:{},
                reporter:{}
            },
            all:{}
        };

        $scope.contribution = function(i, items){
            var value = null;
            var sum = function(ele){
                var s = 0;
                ele.map(function(e){ s += e.rowCount; });
                return s;
            }
            if(items.length != 0){
                value = (parseFloat(i.rowCount) || 0) / (sum(items) * 1.0) * 100;
            }
            return value;
        };

        $scope.widthName = "col-md-7";
        $scope.widthRowCount = "col-md-1";
        $scope.widthProgress = "col-md-4";

        $scope.showAdmin = $scope.isSystemAdministrator();

        var loadData = function(){
            if($scope.isSystemAdministrator()){
                DashboardService.issue.groupAll({ dataType:"assignee"})
                    .then(function(data){
                        $scope.meta.all.assignee = data;
                    })
                    .then(function(){
                        DashboardService.issue.groupAll({ dataType:"reporter"})
                            .then(function(data){
                                $scope.meta.all.reporter = data;
                            })
                    })
                    .then(function(){
                        DashboardService.issue.groupAll({ dataType:"dueDate"})
                            .then(function(data){
                                $scope.meta.all.dueDate = data;
                            })
                    })
                    .then(function(){
                        DashboardService.issue.groupAll({ dataType:"status"})
                            .then(function(data){
                                $scope.meta.all.status = data;
                            })
                    })
                    .then(function(){
                        DashboardService.issue.groupAll({ dataType:"priority"})
                            .then(function(data){
                                $scope.meta.all.priority = data;
                            })
                    })
                    .then(function(){
                        DashboardService.issue.groupAll({ dataType:"project"})
                            .then(function(data){
                                $scope.meta.all.project = data;
                            })
                    });
            } else {
                var id = $scope.currentSession.user.id;
                DashboardService.issue.groupByDataType({ dataType:"project", assignee:id })
                    .then(function(data){
                        $scope.meta.my.assignee.project = data;
                    })
                    .then(function(){
                        DashboardService.issue.groupByDataType({ dataType:"dueDate", assignee:id})
                            .then(function(data){
                                $scope.meta.my.assignee.dueDate = data;
                            })
                    })
                    .then(function(){
                        DashboardService.issue.groupByDataType({ dataType:"status", assignee:id})
                            .then(function(data){
                                $scope.meta.my.assignee.status = data;
                            })
                    })
                    .then(function(){
                        DashboardService.issue.groupByDataType({ dataType:"priority", assignee:id})
                            .then(function(data){
                                $scope.meta.my.assignee.priority = data;
                            })
                    })
                    .then(function(){
                        DashboardService.issue.groupByDataType({ dataType:"project", reporter:id})
                            .then(function(data){
                                $scope.meta.my.reporter.project = data;
                            })
                    })
                    .then(function(){
                        DashboardService.issue.groupByDataType({ dataType:"dueDate", reporter:id})
                            .then(function(data){
                                $scope.meta.my.reporter.dueDate = data;
                            })
                    })
                    .then(function(){
                        DashboardService.issue.groupByDataType({ dataType:"status", reporter:id})
                            .then(function(data){
                                $scope.meta.my.reporter.status = data;
                            })
                    })
                    .then(function(){
                        DashboardService.issue.groupByDataType({ dataType:"priority", reporter:id})
                            .then(function(data){
                                $scope.meta.my.reporter.priority = data;
                            })
                    });
            }

        };

        if($scope.isLoggedIn()){
            loadData();
        }

        $scope.browseIssues = function(i, my, reference){
            var data = { resolved: false },
                user = $scope.currentSession.user.id;
            switch (i.dataType){
                case "project":
                case "priority":
                case "assignee":
                case "reporter":
                    data[i.dataType] = i.id;
                    break;
                case "status":
                    data[i.dataType] = i.name;
                    break;
                case "dueDate":
                    var dueDate = new Date(i.name + '-01');
                    data[i.dataType] = dueDate.toISOString();
                    break;
                default :
                    return;
            }
            if(my){
                if(reference === 'assignee'){
                    data.assignee = user;
                } else {
                    data.reporter = user;
                }
            }

            BrowserService.issue.grid(data);
        }

        $scope.browseProjects = RunApiService.generateUrl(BrowserUrlService.project.grid,{resolved:false});

        $scope.browseProject = function(i){
            BrowserService.project.edit(i.id);
        }
    });
