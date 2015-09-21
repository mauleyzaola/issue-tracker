'use strict';

angular.module("TrackerApp.PermissionScheme.controllers", [])
    .controller("PermissionSchemesController", function($scope, BrowserService, PermissionSchemeService){
        $scope.newItem = function(){
            BrowserService.permissionScheme.add();
        }

        $scope.gridConfig = PermissionSchemeService.gridConfig( { source: PermissionSchemeService.grid });
        $scope.gridParams = {};
    })
    .controller("PermissionSchemeController", function($scope, $routeParams, BrowserService, PermissionSchemeService,
                                                       GroupService, UserService, RoleService, utils, ProjectService){

        $scope.item = { id: $routeParams.id};
        $scope.test = {
            project:null,
            user:null
        };

        var dialogs = {
            projectSelectDialog:$('#projectSelectDialog'),
            userSelectDialog:$('#userSelectDialog')
        };

        $scope.gridProject = ProjectService.gridConfig({
            source: ProjectService.grid,
            columns: [
                { name: "Name", field:"name" },
                { name: "Key", field:"pkey" },
                { name: "Lead", field:"projectLead" },
                { name: "Starts", field:"begins", filter:"dateFormat" },
                { name: "Ends", field:"ends", filter:"dateFormat" },
                { name: "Created", field:"dateCreated", filter:"timeAgo" },
                { name: "Issues", field:"issueCount" },
                { name: "% Completed", field:"percentageCompleted", filter:"percent" }
            ],
            rowClick:function(r){
                dialogs.projectSelectDialog.modal('hide');
                ProjectService.load(r.id)
                    .then(function(data){
                        $scope.test.project = data;
                    });
            }
        });
        $scope.gridProjectParams = {};

        $scope.gridUser = UserService.gridConfig({
            source: UserService.grid,
            rowClick:function(r){
                dialogs.userSelectDialog.modal('hide');
                UserService.load(r.id)
                    .then(function(data){
                        $scope.test.user = data;
                    });
            }
        });
        $scope.gridUserParams = {};



        var resetMetas = function(){
            $scope.meta= {
                group:{
                    selected:[],
                    unselected:[]
                },
                user:{
                    selected:[],
                    unselected:[]
                },
                role:{
                    selected:[],
                    unselected:[]
                }
            };
        };

        var permissionSchemeItems = [];

        var loadPermissionSchemeItems = function(){
            angular.forEach($scope.names, function(name){
                name.items = _.filter(permissionSchemeItems, function(i){ return i.permissionName.id == name.id; });
            });
        };

        var loadItems = function(permissionName){
            PermissionSchemeService.items($routeParams.id)
                .then(function(data){
                    resetMetas();
                    var items = _.filter(data, function(i){ return i.permissionName.id == permissionName.id; }),
                        found = false;

                    angular.forEach($scope.groups, function(g){
                        found = false;
                        angular.forEach(items, function(i){
                            if(i.group){
                                if(i.group.id === g.id){
                                    $scope.meta.group.selected.push(g);
                                    found = true;
                                }
                            }
                        });
                        if(!found){ $scope.meta.group.unselected.push(g); }
                    });

                    angular.forEach($scope.users, function(g){
                        found = false;
                        angular.forEach(items, function(i){
                            if(i.user){
                                if(i.user.id === g.id){
                                    $scope.meta.user.selected.push(g);
                                    found = true;
                                }
                            }
                        });
                        if(!found){ $scope.meta.user.unselected.push(g); }
                    });


                    angular.forEach($scope.roles, function(g){
                        found = false;
                        angular.forEach(items, function(i){
                            if(i.role){
                                if(i.role.id === g.id){
                                    $scope.meta.role.selected.push(g);
                                    found = true;
                                }
                            }
                        });
                        if(!found){ $scope.meta.role.unselected.push(g); }
                    });
                });
        };

        var loadData = function(id){
            PermissionSchemeService.load(id)
                .then(function(data){
                    $scope.item=data;
                })
                .then(function(){
                    PermissionSchemeService.projects(id)
                        .then(function(data){
                            $scope.projects = data;
                        });
                })
                .then(function(){
                    GroupService.list()
                        .then(function(data){
                            $scope.groups = data;
                        })
                        .then(function(){
                            UserService.list()
                                .then(function(data){
                                    $scope.users = data;
                                })
                        })
                        .then(function(){
                            RoleService.list()
                                .then(function(data){
                                    $scope.roles = data;
                                })
                                .then(function(){
                                    PermissionSchemeService.names()
                                        .then(function(data){
                                            return data;
                                        }).then(function(data){
                                            $scope.names = data;
                                            PermissionSchemeService.items(id)
                                                .then(function(data){
                                                    permissionSchemeItems = data;
                                                    loadPermissionSchemeItems();
                                                });
                                        });
                                })
                        });
                });
        }

        if($routeParams.id){
            loadData($routeParams.id);
        }

        $scope.exit = function(){
            BrowserService.permissionScheme.grid();
        }

        $scope.deleteItem = function(){
            if(!utils.confirm()){ return; }
            PermissionSchemeService.remove($scope.item.id).then(function(){ $scope.exit(); });
        }

        $scope.canDelete = function(){
            return $scope.item.id;
        }

        $scope.canSave = function(){
            return true;
        }

        $scope.saveItem = function(){
            var isNewItem = !$scope.item.id;
            PermissionSchemeService.save($scope.item)
                .then(function(data){
                if(!isNewItem){
                    $scope.exit();
                } else {
                    BrowserService.permissionScheme.edit(data.id);
                }
            });
        }

        $scope.editMembers = function(index){
            $scope.selected = $scope.names[index];
            loadItems($scope.selected);
        }

        $scope.addItem = function(item){
            var newItem = {
                permissionScheme:{ id:$scope.item.id },
                permissionName: {id:$scope.selected.id }
            };
            newItem[item.meta.documentType] = item;

            PermissionSchemeService.itemAdd(newItem)
                .then(function(){
                    var items = $scope.meta[item.meta.documentType];
                    var index = utils.findIndexById(item.id, items.unselected);
                    items.unselected.splice(index, 1);
                    items.selected.push(item);

                    permissionSchemeItems.push(newItem);
                    loadPermissionSchemeItems();
                });
        }


        $scope.removeItem = function(item){
            var oldItem = {
                permissionScheme:{ id:$scope.item.id },
                permissionName: {id:$scope.selected.id }
            };
            oldItem[item.meta.documentType] = item;

            PermissionSchemeService.itemRemove(oldItem)
                .then(function(){
                    var items = $scope.meta[item.meta.documentType];
                    var index = utils.findIndexById(item.id, items.selected);
                    items.selected.splice(index, 1);
                    items.unselected.push(item);

                    index = -1;

                    for(var i=0; i<permissionSchemeItems.length; i++){
                        var ps = permissionSchemeItems[i];
                        if( (oldItem.group && ps.group && oldItem.group.id == ps.group.id
                                && oldItem.permissionName.id == ps.permissionName.id) ||
                            (oldItem.role && ps.role && oldItem.role.id == ps.role.id
                            && oldItem.permissionName.id == ps.permissionName.id) ||
                            (oldItem.user&& ps.user && oldItem.user.id == ps.user.id
                            && oldItem.permissionName.id == ps.permissionName.id)){
                            index = i;
                            break;
                        }
                    }

                    if(index != -1){
                        permissionSchemeItems.splice(index, 1);
                        loadPermissionSchemeItems();
                    }
                });
        }

        $scope.parsePermissionSchemeItem = function(item){
            if(item.group){
                return item.group.name;
            } else if (item.user) {
                return item.user.name + ' ' + item.user.lastName;
            } else if (item.role){
                return item.role.name;
            }
        }

        $scope.testPermissions = function(){
            PermissionSchemeService.permissionAvailableUser({
                issue:{project: $scope.test.project },
                user: $scope.test.user
            }).then(function(data){
                var items = [];
                angular.forEach($scope.names, function(name){
                    var exists = (_.find(data, function(i){ return name.id === i.id; }) != null );
                    items.push({ name: name.meta.friendName, exists: exists});
                });
                $scope.testPermissionItems = items;
            });
        }

        $scope.clearAll = function(){
            if(!utils.confirm("Are you sure to remove all the current permissions?")){ return; }
            PermissionSchemeService.clear($scope.item.id)
                .then(function(){
                    loadData($routeParams.id);
                })
        }
    });
