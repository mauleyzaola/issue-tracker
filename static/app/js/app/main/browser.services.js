'use strict';

angular.module("TrackerApp.Browser.services", [])

    .factory("BrowserUrlService", function(){
        return{

            group:{
                add:"/catalog/group",
                grid:"/catalog/groups",
                edit:function(id){ return "/catalog/group/" + id; }
            },

            role:{
                add:"/catalog/role",
                grid:"/catalog/roles",
                edit:function(id){ return "/catalog/role/" + id; }
            },

            user:{
                add:"/catalog/user",
                grid:"/catalog/users",
                edit:function(id){ return "/catalog/user/" + id; },
                myProfile:"/account/myprofile"
            },

            issue:{
                add:"/issue/new",
                grid:"/issue/issues/list",
                edit:function(pkey){ return "/issue/browse/" + pkey; }
            },

            issueAttachment:{
                edit:function(id){ return '/issue/attachment/' + id ; }
            },

            permissionScheme:{
                add:"/catalog/permissionscheme",
                grid:"/catalog/permissionschemes",
                edit:function(id){ return "/catalog/permissionscheme/" + id; }
            },

            priority:{
                add:"/catalog/priority",
                grid:"/catalog/priorities",
                edit:function(id){ return "/catalog/priority/" + id; }
            },

            project:{
                add:"/issue/project",
                grid:"/issue/projects",
                edit:function(id){ return "/issue/project/" + id; }
            },

            workflow:{
                add:"/catalog/workflow",
                grid:"/catalog/workflows",
                edit:function(id){ return "/catalog/workflow/" + id; }
            },

            workflowStep:{

            }

        }
    })
    .factory("BrowserService", function($location, BrowserUrlService){
        var clearQueryString = function(search){
            if(!search){
                $location.$$search = {};
            }else{
                $location.$$search = search;
            }
        };

        return {

            group:{
                add:function(search){
                    clearQueryString(search);
                    $location.path(BrowserUrlService.group.add);
                },
                grid:function(search){
                    clearQueryString(search);
                    $location.path(BrowserUrlService.group.grid);
                },
                edit:function(id){
                    clearQueryString();
                    $location.path(BrowserUrlService.group.edit(id));
                }
            },

            role:{
                add:function(search){
                    clearQueryString(search);
                    $location.path(BrowserUrlService.role.add);
                },
                grid:function(search){
                    clearQueryString(search);
                    $location.path(BrowserUrlService.role.grid);
                },
                edit:function(id){
                    clearQueryString();
                    $location.path(BrowserUrlService.role.edit(id));
                }
            },

            user:{
                add:function(search){
                    clearQueryString(search);
                    $location.path(BrowserUrlService.user.add);
                },
                grid:function(search){
                    clearQueryString(search);
                    $location.path(BrowserUrlService.user.grid);
                },
                edit:function(id){
                    clearQueryString();
                    $location.path(BrowserUrlService.user.edit(id));
                }
            },

            issue:{
                add:function(search){
                    clearQueryString(search);
                    $location.path(BrowserUrlService.issue.add);
                },
                grid:function(search){
                    clearQueryString(search);
                    $location.path(BrowserUrlService.issue.grid);
                },
                edit:function(pkey, queryParams){
                    clearQueryString(queryParams);
                    $location.path(BrowserUrlService.issue.edit(pkey));
                }
            },

            issueAttachment:{
                edit:function(id){
                    clearQueryString();
                    $location.path(BrowserUrlService.issueAttachment.edit(id));
                }
            },

            permissionScheme:{
                add:function(search){
                    clearQueryString(search);
                    $location.path(BrowserUrlService.permissionScheme.add);
                },
                grid:function(search){
                    clearQueryString(search);
                    $location.path(BrowserUrlService.permissionScheme.grid);
                },
                edit:function(id){
                    clearQueryString();
                    $location.path(BrowserUrlService.permissionScheme.edit(id));
                }
            },

            priority:{
                add:function(search){
                    clearQueryString(search);
                    $location.path(BrowserUrlService.priority.add);
                },
                grid:function(search){
                    clearQueryString(search);
                    $location.path(BrowserUrlService.priority.grid);
                },
                edit:function(id){
                    clearQueryString();
                    $location.path(BrowserUrlService.priority.edit(id));
                }
            },

            project:{
                add:function(search){
                    clearQueryString(search);
                    $location.path(BrowserUrlService.project.add);
                },
                grid:function(search){
                    clearQueryString(search);
                    $location.path(BrowserUrlService.project.grid);
                },
                edit:function(id){
                    $location.path(BrowserUrlService.project.edit(id));
                }
            },

            workflow:{
                add:function(search){
                    clearQueryString(search);
                    $location.path(BrowserUrlService.workflow.add);
                },
                grid:function(search){
                    clearQueryString(search);
                    $location.path(BrowserUrlService.workflow.grid);
                },
                edit:function(id){
                    clearQueryString();
                    $location.path(BrowserUrlService.workflow.edit(id));
                }
            },

            workflowStep:{

            }
        }
    })

    .factory("ResolveUrlService", function($log, BrowserUrlService, NotificationTypes){
        return {
            resolveUrl:function(item){
                if(!item){
                    return;
                }
                if (!item.meta){
                    return;
                }

                try{
                    var objectType = item.meta.documentType;
                    if(!objectType || objectType.length === 0){ return; }
                    var identifier;
                    if(objectType === NotificationTypes.objectType.issue){
                        identifier = item.pkey;
                    } else if(objectType === NotificationTypes.objectType.bomItemFile) {
                        identifier = item.bomItem.id;
                    } else {
                        identifier = item.id;
                    }

                    return BrowserUrlService[objectType].edit(identifier);
                } catch(e) {
                    $log.warn(e);
                    return null;
                }
            }
        }
    })

