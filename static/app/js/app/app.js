'use strict';


angular.module('TrackerApp', [
    'TrackerApp.controllers', 'TrackerApp.filters', 'TrackerApp.services', 'TrackerApp.Directives',
    'TrackerApp.ApiPath.services', 'TrackerApp.Notification.services', 'TrackerApp.QueryString.services','TrackerApp.Browser.services',

    'TrackerApp.Account.controllers', 'TrackerApp.Account.services',
    'TrackerApp.Group.controllers','TrackerApp.Group.services',
    'TrackerApp.Role.controllers','TrackerApp.Role.services',
    'TrackerApp.User.controllers', 'TrackerApp.User.services',
    'TrackerApp.Issue.controllers','TrackerApp.Issue.services',
    'TrackerApp.PermissionScheme.controllers','TrackerApp.PermissionScheme.services',
    'TrackerApp.Priority.controllers','TrackerApp.Priority.services',
    'TrackerApp.Project.controllers','TrackerApp.Project.services',
    'TrackerApp.Workflow.controllers', 'TrackerApp.Workflow.services',

    'TrackerApp.Account.Directives',
    'TrackerApp.Issue.Directives',
    'TrackerApp.Catalog.services',
    'TrackerApp.Grid.Directives',
    'TrackerApp.FileItem.services',

    'TrackerApp.Dashboard.services',


    'ngFileUpload',
    'ui.bootstrap',
    'ngRoute',
    'ngSanitize'
])
    .config(function ($routeProvider, $locationProvider, $httpProvider, $provide ) {

    //Avoid caching issues with IE
    if (!$httpProvider.defaults.headers.get) {
        $httpProvider.defaults.headers.get = {};
    }
    //disable IE ajax request caching
    $httpProvider.defaults.headers.get['If-Modified-Since'] = '0';


    var genericGridTemplatePath = '/templates/generic/generic-grid.html';

    $routeProvider
        .when ('/', { templateUrl: '/templates/main.html', controller: 'Index.Controller' })

        .when('/account/login',{ templateUrl: '/templates/account/login.html', controller: 'Login.Controller' })
        .when('/account/logout',{ templateUrl: '/templates/account/logout.html', controller: 'Logout.Controller' })
        .when('/account/myprofile',{ templateUrl: '/templates/account/myprofile.html', controller: 'MyProfile.Controller' })
        .when('/account/passwordtoken/:token', { templateUrl: '/templates/account/passwordtoken.html', controller: 'Login.RecoverPassword' })

        .when('/catalog/groups', { templateUrl: genericGridTemplatePath, controller: 'GroupsController' })
        .when('/catalog/group/:id', { templateUrl: '/templates/catalog/group/group.html', controller: 'GroupController' })
        .when('/catalog/group', { templateUrl: '/templates/catalog/group/group.html', controller: 'GroupController' })

        .when('/catalog/permissionschemes', { templateUrl: genericGridTemplatePath, controller: 'PermissionSchemesController' })
        .when('/catalog/permissionscheme/:id', { templateUrl: '/templates/catalog/permission_scheme/permission_scheme.html', controller: 'PermissionSchemeController' })
        .when('/catalog/permissionscheme', { templateUrl: '/templates/catalog/permission_scheme/permission_scheme.html', controller: 'PermissionSchemeController' })

        .when('/catalog/roles', { templateUrl: genericGridTemplatePath, controller: 'RolesController' })
        .when('/catalog/role/:id', { templateUrl: '/templates/catalog/role/role.html', controller: 'RoleController' })
        .when('/catalog/role', { templateUrl: '/templates/catalog/role/role.html', controller: 'RoleController' })

        .when('/catalog/users', { templateUrl: genericGridTemplatePath, controller: 'UsersController' })
        .when('/catalog/user/:id', { templateUrl: '/templates/catalog/user/user.html', controller: 'UserController' })
        .when('/catalog/user', { templateUrl: '/templates/catalog/user/user.html', controller: 'UserController' })

        .when('/catalog/workflows', { templateUrl: genericGridTemplatePath, controller: 'WorkflowsController' })
        .when('/catalog/workflow/:id', { templateUrl: '/templates/issue/workflow/workflow.html', controller: 'WorkflowController' })
        .when('/catalog/workflow', { templateUrl: '/templates/issue/workflow/workflow.html', controller: 'WorkflowController' })

        .when('/catalog/priorities', { templateUrl: genericGridTemplatePath, controller: 'PrioritiesController' })
        .when('/catalog/priority/:id', { templateUrl: '/templates/catalog/priority/priority.html', controller: 'PriorityController' })
        .when('/catalog/priority', { templateUrl: '/templates/catalog/priority/priority.html', controller: 'PriorityController' })

        .when('/issue/issues/list', { templateUrl: genericGridTemplatePath, controller: 'IssuesController' })
        .when('/issue/browse/:pkey', { templateUrl: '/templates/issue/issue/issue.html', controller: 'IssueController' })
        .when('/issue/new', { templateUrl: '/templates/issue/issue/issue.html', controller: 'IssueController' })

        .when('/issue/attachment/:id', { templateUrl: genericGridTemplatePath, controller: 'IssueAttachmentForwarderController' })

        .when('/issue/projects', { templateUrl: '/templates/issue/project/mosaic.html', controller: 'ProjectsMosaic' })
        .when('/issue/project/:id', { templateUrl: '/templates/issue/project/project.html', controller: 'ProjectController' })
        .when('/issue/project', { templateUrl: '/templates/issue/project/project.html', controller: 'ProjectController' })


        //.otherwise({ redirectTo: '/' });


    $locationProvider.html5Mode({
        enabled:true,
        requireBase:false
    });

        /**
         * Intercept http calls.
         */
    $provide.factory('MyHttpInterceptor', function ($q, $location, $rootScope, $timeout, $window, $log, Notifier, SessionManagement) {
        return {
            request: function (config) {
                var session = SessionManagement.currentSession();
                config.headers["accept"] = "application/json";
                if(session && session.id){
                    config.headers["token"] = session.id;
                }
                return config || $q.when(config);
            },
            requestError: function (rejection) { return $q.reject(rejection); },
            response: function (response) { return response || $q.when(response); },
            responseError: function (rejection) {
                $window.console.log(rejection);
                var loginUrl = '/account/login';
                var msg = rejection.data.Message || 'Error from API';

                switch (rejection.status){
                    case 401:
                        if($location.path() == loginUrl){ Notifier.error({title:'Error', text: msg});}
                        else {
                            $location.$$search = {returnurl:  $location.url()};
                            $location.path(loginUrl);
                        }
                        break;
                    case 404:
                        Notifier.error({title:'not found', text: rejection.data }); break;
                    default :
                        Notifier.warning({title:'API', text: msg});
                        break;
                }

                // Return the promise rejection.
                return $q.reject(rejection);
            }
        };
    });

    // Add the interceptor to the $httpProvider.
    $httpProvider.interceptors.push('MyHttpInterceptor');
});
