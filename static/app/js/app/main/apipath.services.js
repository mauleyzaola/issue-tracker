'use strict';

angular.module("TrackerApp.ApiPath.services", [])
    .factory("PathResolver", function(){
        return {
            baseApiAuthPath: "/api/auth",
            baseApiNoAuthPath: "/api/noauth",
            baseApiSysAdminPath: "/api/admin",
            baseFileDownload: "/files/"
        }
    })

    .factory('PathService', function(PathResolver, QueryStringNames, utils){
        return {
            account:{
                login: PathResolver.baseApiNoAuthPath + "/account/login",
                logout : PathResolver.baseApiAuthPath + "/account/logout",
                profile: PathResolver.baseApiAuthPath + "/account/myprofile",
                removeToken: function(token){ return PathResolver.baseApiAuthPath + "/account/session/" + token; },
                getTokens: PathResolver.baseApiAuthPath + "/account/token/list",
                passwordRecoverToken: PathResolver.baseApiNoAuthPath + "/account/passwordrecover",
                verifyToken: PathResolver.baseApiNoAuthPath + "/account/verifytoken",
                changeMyPassword: PathResolver.baseApiAuthPath + "/account/changemypassword",
                changePassword: PathResolver.baseApiSysAdminPath + "/account/changepassword"
            },

            file:{
                download:PathResolver.baseApiAuthPath + "/file",
                upload:PathResolver.baseApiAuthPath + "/file",
                directoryGrid:PathResolver.baseApiAuthPath + '/files/directory/grid',
                fileGrid:PathResolver.baseApiAuthPath + '/files/file/grid'
            },

            group: {
                load: function(id) { return PathResolver.baseApiAuthPath + "/group/" + id; },
                save: PathResolver.baseApiSysAdminPath + "/group",
                remove: function(id) { return PathResolver.baseApiSysAdminPath + "/group/" + id; },
                grid:PathResolver.baseApiAuthPath + "/group/grid",
                list: PathResolver.baseApiAuthPath + "/groups",
                groups:function(id){ return PathResolver.baseApiAuthPath + "/user/" + id + "/groups" },
                users:function(id){ return PathResolver.baseApiAuthPath + "/group/" + id + "/users" },
                addGroupUser: PathResolver.baseApiSysAdminPath + "/group/users/add",
                removeGroupUser: PathResolver.baseApiSysAdminPath + "/group/users/remove"
            },

            issue: {
                createMeta: function(data){
                    var url = PathResolver.baseApiAuthPath + "/issue/createmeta";
                    if(data.pkey){
                        url += "/" + data.pkey;
                    } else {
                        data.pkey = "";
                    }
                    url += "?" + utils.queryStringFromObject(data);
                    return url;
                },
                load: function(data){
                    if(data.id){
                        return PathResolver.baseApiAuthPath + "/issue/" + data.id;
                    } else {
                        return PathResolver.baseApiAuthPath + "/issue?" + utils.queryStringFromObject(data);
                    }
                },
                getChildren:function(id){return PathResolver.baseApiAuthPath + "/issue/" + id + "/children"},
                move:function(id){ return PathResolver.baseApiAuthPath + "/issue/" + id + "/move"; },
                save: function (data) {
                    var parent = '';
                    if(data.parent){
                        parent=data.parent.id;
                    }
                    return PathResolver.baseApiAuthPath + "/issue?" + QueryStringNames.parent + "=" + parent;
                },
                remove: function(id){ return PathResolver.baseApiAuthPath + "/issue/" + id; },
                grid: PathResolver.baseApiAuthPath + "/issue/grid",

                assignToMe:function(data){ return PathResolver.baseApiAuthPath + "/issue/assigntome?" + utils.queryStringFromObject(data); },
                reporterIsMe:function(data){ return PathResolver.baseApiAuthPath + "/issue/reporterisme?" + utils.queryStringFromObject(data); },

                changeStatus: function(data){ return PathResolver.baseApiAuthPath + "/issue/status?" + utils.queryStringFromObject(data); },

                commentAdd:function(data){ return PathResolver.baseApiAuthPath + "/issue/" + data.issue.id + "/comment" },
                commentUpdate: function(data){ return PathResolver.baseApiAuthPath + "/issue/" + data.issue.id + "/comment" },
                commentRemove: function(data){ return PathResolver.baseApiAuthPath + "/issue/" + data.id + "/comment" },
                commentList: function(data){ return PathResolver.baseApiAuthPath + "/issue/" + data.id + "/comments"},

                subscriptionToggle:function(data){ return PathResolver.baseApiAuthPath + "/issue/" + data.id + "/subscription/toggle" },
                mySubscriptions: PathResolver.baseApiAuthPath + "/issue/mysubscriptions",
                subscribedSelected:function(id){ return PathResolver.baseApiAuthPath + '/issue/' +  id + '/subscribedselected' ;},
                subscriptionToggleAny: PathResolver.baseApiAuthPath + '/issue/subscription/toggle/any',

                attachments:function(data){ return PathResolver.baseApiAuthPath + "/issue/" + data.id + "/attachments"; },
                attachmentAdd:function(data){ return PathResolver.baseApiAuthPath + "/issue/" + data.issue.id + "/attachment";},
                attachmentRemove:function(data){ return PathResolver.baseApiAuthPath + "/issue/" + data.id + "/attachment"; },
                attachmentLoad:function(id){ return PathResolver.baseApiAuthPath + '/issue/attachment/' + id; },
                currentUserSubscribed:function(id){ return PathResolver.baseApiAuthPath + "/issue/" + id + "/subscribed"; },

                /* dashboard and reporting */
                groupAll:function(params){ return PathResolver.baseApiAuthPath + "/issue/groupall?" + utils.queryStringFromObject(params); },
                groupByDataType:PathResolver.baseApiAuthPath + '/issue/group/bydatatype'
            },

            permissionScheme: {
                load: function(id) { return PathResolver.baseApiAuthPath + "/permissionscheme/" + id; },
                save: PathResolver.baseApiSysAdminPath + "/permissionscheme",
                remove: function(id) { return PathResolver.baseApiSysAdminPath + "/permissionscheme/" + id; },
                clear: function(id) { return PathResolver.baseApiAuthPath + "/permissionscheme/" + id + "/clear"; },
                grid: PathResolver.baseApiAuthPath + "/permissionscheme/grid",
                list: PathResolver.baseApiAuthPath + "/permissionschemes",
                names:PathResolver.baseApiAuthPath + "/permissionnames",
                projects:function(id){ return PathResolver.baseApiAuthPath + "/permissionscheme/" + id + "/projects"; }
            },

            permissionSchemeItem: {
                add: PathResolver.baseApiSysAdminPath + "/permissionschemeitem/add",
                remove: PathResolver.baseApiSysAdminPath + "/permissionschemeitem/remove",
                list: function(id){ return PathResolver.baseApiAuthPath + "/permissionschemeitems/" + id;},
                permissionAvailableUser:PathResolver.baseApiAuthPath + "/permissionschemeitem/user/available"
            },

            priority: {
                load: function(id) { return PathResolver.baseApiAuthPath + "/priority/" + id; },
                save: PathResolver.baseApiAuthPath + "/priority",
                remove: function(id) { return PathResolver.baseApiAuthPath + "/priority/" + id; },
                grid: PathResolver.baseApiAuthPath + "/priority/grid",
                list: PathResolver.baseApiAuthPath + "/priorities"
            },

            project: {
                createMeta: function(id){ return PathResolver.baseApiAuthPath + "/project/createmeta?" + QueryStringNames.id + "=" + id; },
                load: function(id){ return PathResolver.baseApiAuthPath + "/project/" + id; },
                create: PathResolver.baseApiAuthPath + "/project",
                update: PathResolver.baseApiAuthPath + "/project",
                remove: function(id) { return PathResolver.baseApiAuthPath + "/project/" + id; },
                grid: PathResolver.baseApiAuthPath + "/project/grid",
                projectRoles:function(id){
                    return PathResolver.baseApiAuthPath + "/projectrole/project/members?" + QueryStringNames.project + "=" + (id || "");
                },
                projectRoleMembers:function(id){ return PathResolver.baseApiAuthPath + "/projectrole/" + id + "/members"; },
                projectRoleProjectMembers:function(id){ return PathResolver.baseApiAuthPath + '/project/' + id + '/members';},
                projectRoleMemberAdd:PathResolver.baseApiAuthPath + "/projectrole/members/add",
                projectRoleMemberRemove:PathResolver.baseApiAuthPath + "/projectrole/members/remove"
            },

            role: {
                load: function(id){ return PathResolver.baseApiAuthPath + "/role/" + id; },
                save: PathResolver.baseApiSysAdminPath + "/role",
                remove: function(id) { return PathResolver.baseApiSysAdminPath + "/role/" + id; },
                grid: PathResolver.baseApiAuthPath + "/role/grid",
                list: PathResolver.baseApiAuthPath + "/roles"
            },

            status:{
                load: function(id){ return PathResolver.baseApiAuthPath + "/status/" + id; },
                save: PathResolver.baseApiSysAdminPath + "/status",
                remove: function(id){ return PathResolver.baseApiSysAdminPath + "/status/" + id; },
                list: function(workflow){ return PathResolver.baseApiAuthPath + "/statuses/" + workflow; }
            },

            user: {
                load: function(id){ return PathResolver.baseApiAuthPath + "/user/" + id; },
                changePassword: PathResolver.baseApiSysAdminPath + "/user/changepassword",
                save: PathResolver.baseApiSysAdminPath + "/user",
                remove: function(id) { return PathResolver.baseApiSysAdminPath + "/user/" + id; },
                grid: PathResolver.baseApiAuthPath + "/user/grid",
                list: PathResolver.baseApiAuthPath + "/users"
            },

            workflow: {
                load: function(id){ return PathResolver.baseApiAuthPath + "/workflow/" + id; },
                createMeta: function(id){ return PathResolver.baseApiAuthPath + '/workflow/' + id + '/createmeta'; },
                save: PathResolver.baseApiSysAdminPath + "/workflow",
                remove: function(id) { return PathResolver.baseApiSysAdminPath + "/workflow/" + id; },
                grid: PathResolver.baseApiAuthPath + "/workflow/grid",
                list: PathResolver.baseApiAuthPath + "/workflows"
            },

            workflowStep: {
                save: PathResolver.baseApiSysAdminPath + "/workflowstep",
                remove: function(id){ return PathResolver.baseApiSysAdminPath + "/workflowstep/" + id; },
                list: function(workflow){ return PathResolver.baseApiAuthPath + "/workflowsteps/" + workflow; },
                availableSteps: function(data){ return PathResolver.baseApiAuthPath + "/workflowsteps/available/" + data.workflow + "?" + QueryStringNames.status + "=" + data.status; },
                availableStepsUser: function(data){ return PathResolver.baseApiAuthPath + "/workflowsteps/user/" + data.workflow + "?" + QueryStringNames.status + "=" + data.status; },
                members:function(id){ return PathResolver.baseApiAuthPath + "/workflowstepmembers/" + id; },
                memberAdd:PathResolver.baseApiSysAdminPath + "/workflowstepmember/add",
                memberRemove:PathResolver.baseApiSysAdminPath + "/workflowstepmember/remove",
                memberGroups:function(id){ return PathResolver.baseApiSysAdminPath + "/workflowstep/" + id + "/groups"; },
                memberUsers:function(id){ return PathResolver.baseApiSysAdminPath + "/workflowstep/" + id + "/users"; }
            }
        }
    })

    .factory('RunApiService', function(utils){
        return {
            generateUrl:function(baseUrl, params){
                if(!params){
                    return baseUrl;
                } else {
                    return baseUrl + "?" + utils.queryStringFromObject(params);
                }
            }
        }
    })
