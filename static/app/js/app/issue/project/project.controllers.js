'use strict';

angular.module("TrackerApp.Project.controllers", [])
    .controller("ProjectsMosaic", function($scope, $location, BrowserService, ProjectService, $filter, dashboardFn){
        $scope.newItem = function(){
            BrowserService.project.add();
        };

        var pars = $location.$$search;
        angular.extend(pars, {
            sortField:'name',
            sortDirection:'asc'
        });


        $scope.getCss = function(i){
            var value = dashboardFn.cssColor(i.percentageCompleted);
            if(value === 'default'){
                value = 'info';
            }
            return 'alert alert-' + value;
        }

        $scope.getType = function(i){
            var value = dashboardFn.cssColor(i.percentageCompleted);
            if(value === 'default'){
                value = 'info';
            }
            return value;
        }

        $scope.browseItem = function(i){
            $location.$$search = { resolved: $location.$$search.resolved };
            BrowserService.project.edit(i.id);
        }


        var paginationPars = {
            pageNumber:0,
            pageSize:10
        };

        $scope.pagination = {
            currentPage:1,
            maxSize:6,
            pageCount:0,
            pageSize:10,
            totalCount:0
        };

        var loadData = function(){
            paginationPars.pageNumber = $scope.pagination.currentPage || 1;
            paginationPars.size = $scope.pagination.pageSize;
            var params = paginationPars;
            angular.extend(params, pars);
            ProjectService.grid(params).then(function(data){
                var pageCount = parseInt(data.totalCount / $scope.pagination.pageSize),
                    remainder = data.totalCount % $scope.pagination.pageSize;
                pageCount += Math.sign(remainder);
                $scope.pagination.pageCount = pageCount;
                $scope.pagination.totalCount = data.totalCount;
                $scope.items = data.rows;
            });
        };

        $scope.$watch('pagination.currentPage', function(newValue){
            if(!newValue){ return; }
                loadData();
        });
    })
    .controller("ProjectsController", function($scope, BrowserService, ProjectService){
        $scope.newItem = function(){
            BrowserService.project.add();
        };
        $scope.gridConfig = ProjectService.gridConfig({ source: ProjectService.grid});
        $scope.gridParams = {};
    })
    .controller("ProjectController", function($scope, $location, $routeParams,
                                              BrowserService, $timeout, ProjectService, IssueService, utils, QueryStringNames,
                                              PermissionSchemeService, GroupService, UserService,
                                              DashboardService){
        $scope.item = { id: $routeParams.id};

        $scope.dashboard = {
            all:{}
        };
        $scope.widthName = "col-md-7";
        $scope.widthRowCount = "col-md-1";
        $scope.widthProgress = "col-md-4";

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



        $scope.browseIssues = function(i, my, reference){
            var data = {
                resolved: false,
                project: $scope.item.id
            };
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

            BrowserService.issue.grid(data);
        }


        var projectRoleMembers = [];

        var resetMetas = function(){
            $scope.meta= {
                group:{
                    selected:[],
                    unselected:[]
                },
                user:{
                    selected:[],
                    unselected:[]
                }
            };
        };

        var loadProjectRoleMembers = function(){
            angular.forEach($scope.projectRoles, function(item){
                item.items = _.filter(projectRoleMembers, function(i){ return i.projectRole.id == item.id; });
            });
        };

        var loadMembers = function(projectRole){
            ProjectService.projectRoleMembers(projectRole.id)
                .then(function(data){
                    resetMetas();
                    var found = false;

                    angular.forEach($scope.groups, function(g){
                        found = false;
                        angular.forEach(_.filter(data, function(i){ return i.group; }), function(i){
                            if(i.group.id === g.id){
                                $scope.meta.group.selected.push(g);
                                found = true;
                            }
                        });
                        if(!found){ $scope.meta.group.unselected.push(g); }
                    });

                    angular.forEach($scope.users, function(g){
                        found = false;
                        angular.forEach(_.filter(data, function(i){ return i.user; }), function(i){
                            if(i.user.id === g.id){
                                $scope.meta.user.selected.push(g);
                                found = true;
                            }
                        });
                        if(!found){ $scope.meta.user.unselected.push(g); }
                    });
                });
        };


        ProjectService.createMeta($routeParams.id)
            .then(function(data){
                $scope.item = data.item;
                $scope.users = data.users;
                $scope.projectRoles = data.projectRoles;
            }).then(function(){
                PermissionSchemeService.list()
                    .then(function(data){
                        data.unshift({name:"None"});
                        $scope.permissionSchemes = data;
                    });
            }).then(function(){
                if($scope.item.id){
                    ProjectService.projectRoleProjectMembers($routeParams.id)
                        .then(function(data){
                            projectRoleMembers = data;
                        }).then(function(){
                            GroupService.list()
                                .then(function(data){
                                    $scope.groups = data;
                                })
                                .then(function(){
                                    UserService.list()
                                        .then(function(data){
                                            $scope.users = data;
                                            loadProjectRoleMembers();
                                        });
                                });
                        });
                }
            });

        if($routeParams.id){
            DashboardService.issue.groupByDataType({ dataType:"assignee", project:$routeParams.id})
                .then(function(data){
                    $scope.dashboard.all.assignee = data;
                })
                .then(function(){
                    DashboardService.issue.groupByDataType({ dataType:"reporter", project:$routeParams.id})
                        .then(function(data){
                            $scope.dashboard.all.reporter = data;
                        })
                })
                .then(function(){
                    DashboardService.issue.groupByDataType({ dataType:"status", project:$routeParams.id})
                        .then(function(data){
                            $scope.dashboard.all.status = data;
                        })
                })
                .then(function(){
                    DashboardService.issue.groupByDataType({ dataType:"priority", project:$routeParams.id})
                        .then(function(data){
                            $scope.dashboard.all.priority = data;
                        })
                })
                .then(function(){
                    DashboardService.issue.groupByDataType({ dataType:"dueDate", project:$routeParams.id})
                        .then(function(data){
                            $scope.dashboard.all.dueDate = data;
                        })
                })
        }

        $scope.exit = function(){
            BrowserService.project.grid($location.$$search);
        };

        $scope.deleteItem = function(){
            if(!utils.confirm()){ return; }
            ProjectService.remove($scope.item.id).then(function(){ $scope.exit(); });
        };

        $scope.canDelete = function(){
            return $scope.item.id;
        };

        $scope.canSave = function(){
            return true;
        };

        $scope.saveItem = function(){
            var isNewItem = !$scope.item.id;
            ProjectService.save($scope.item).then(function(data){
                if(isNewItem){
                    BrowserService.project.edit(data.id);
                } else {
                    $scope.exit();
                }
            });
        };

        if($routeParams.id){
            $scope.gridIssue= IssueService.gridConfig({
                source: IssueService.grid,
                columns: [
                    { name: "Key", field:"pkey" },
                    { name: "Name", field:"name" },
                    { name: "Assignee", field:"assignee" },
                    { name: "Reporter", field:"reporter" },
                    { name: "Priority", field:"priority" },
                    { name: "Status", field:"status" },
                    { name: "Due Date", field:"dueDate", filter:"timeAgo" }
                ]
            });
            $scope.gridIssueParams = { project: $routeParams.id };
        }

        $scope.addIssue = function(){
            $location.$$search[QueryStringNames.project]=$scope.item.id;
            $location.path("/issue/new");
        }

        $scope.editMembers = function(index){
            $scope.selected = $scope.projectRoles[index];
            loadMembers($scope.selected);
        }

        $scope.addMember = function(item){
            var newItem = {
                projectRole: {
                    project:{ id:$scope.item.id },
                    role: { id: $scope.selected.role.id },
                    id: $scope.selected.id
                }
            };
            newItem[item.meta.documentType] = item;

            ProjectService.projectRoleMemberAdd(newItem)
                .then(function(){
                    var items = $scope.meta[item.meta.documentType];
                    var index = utils.findIndexById(item.id, items.unselected);
                    items.unselected.splice(index, 1);
                    items.selected.push(item);
                    projectRoleMembers.push(newItem);
                    loadProjectRoleMembers()
                });
        }


        $scope.removeMember = function(item){
            var oldItem = {
                projectRole: {
                    project:{ id:$scope.item.id },
                    role: { id: $scope.selected.role.id },
                    id: $scope.selected.id
                }
            };
            oldItem[item.meta.documentType] = item;

            ProjectService.projectRoleMemberRemove(oldItem)
                .then(function(){
                    var items = $scope.meta[item.meta.documentType];
                    var index = utils.findIndexById(item.id, items.selected);
                    items.selected.splice(index, 1);
                    items.unselected.push(item);

                    index = -1;

                    for(var i=0; i<projectRoleMembers.length; i++){
                        var ps = projectRoleMembers[i];
                        if( (oldItem.group && ps.group && oldItem.group.id == ps.group.id
                            && oldItem.projectRole.id == ps.projectRole.id) ||
                            (oldItem.user&& ps.user && oldItem.user.id == ps.user.id
                            && oldItem.projectRole.id == ps.projectRole.id)){
                            index = i;
                            break;
                        }
                    }

                    if(index != -1){
                        projectRoleMembers.splice(index, 1);
                        loadProjectRoleMembers()
                    }
                });
        }

        $scope.parseProjectRoleMember = function(item){
            if(item.group){
                return item.group.name;
            } else if (item.user) {
                return item.user.name + ' ' + item.user.lastName;
            }
        }


    });
