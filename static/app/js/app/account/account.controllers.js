'use strict';

angular.module("TrackerApp.Account.controllers", [])
    .controller("Login.RecoverPassword", function($rootScope, $scope, $location, $window, $routeParams, $timeout, SessionManagement, AccountService, Notifier){

        AccountService.verifyToken($routeParams.token)
            .then(function(data){
                Notifier.push({
                    title:"Token validation was successful",
                    text:"Need to reset your password"
                });

                SessionManagement.currentSession(data);
                $rootScope.currentSession = data;

                $scope.$emit("user:login", {} );

                $timeout(function(){
                    $location.$$search = { tab: 3};
                    $location.path("/account/myprofile");
                }, 300);
            });
    })
    .controller("Login.Controller", function($http, $rootScope, $scope, $location, AccountService, SessionManagement){
        $scope.userInfo = {
            email: $location.$$search.email,
            password: null
        }


        SessionManagement.clearCurrentSession();

        $scope.exit = function(){
            $location.search("");
            $location.path("/");
        }

        $scope.login = function(){
            AccountService.login({email: $scope.userInfo.email, password: $scope.userInfo.password})
                .then(function(data){
                    $scope.$emit('user:login', data);
                })
        };

        $scope.submitRecoverPassword = function(){
            AccountService.passwordRecoverToken($scope.userInfo.email)
                .then(function(){
                    $location.$$search={};
                    $location.path("/");
                });
        };
    })
    .controller("Logout.Controller", function($rootScope, $scope, $location, $timeout, AccountService){
        AccountService.logout()
            .then(function(){
                $scope.$emit('user:logout', {});
            });

    })
    .controller("MyProfile.Controller", function($scope, $rootScope, $location, $window, AccountService, BusinessUnitService,
                                                 SessionManagement, utils, IssueService, BrowserUrlService){

        $scope.selectedTab = $location.$$search.tab || 1;
        $scope.sesiones = [];
        $scope.item={};
        $scope.issueSubscriptions = [];

        AccountService.loadProfile()
            .then(function(data){
                $scope.item = data;
            })
            .then(function(){
                AccountService.getTokens()
                    .then(function(data){
                        $scope.sessions=data;
                    });
            });

        $scope.deleteToken = function(index, item){
            if(!utils.confirm()) { return; }
            AccountService.removeToken(item.id)
                .then(function(){
                    $scope.sessions.splice(index,1);
                });
        }

        /**
         * Returns the warning password text depending on the input from the user
         * @method passwordText
         * @returns {string}
         */
        $scope.passwordText = function(){
            if(!$scope.item){
                return null;
            } else if(!$scope.password1 || $scope.password1.length == 0){
                return "New password cannot be empty";
            } else if($scope.password1 != $scope.password2){
                return "Passwords don't match"
            } else {
                return null;
            }
        }


        $scope.$watch('password1', function (newVal, oldVal) {
            if (newVal !== oldVal) {
                $scope.passwordText();
            }
        }, true);
        $scope.$watch('password2', function (newVal, oldVal) {
            if (newVal !== oldVal) {
                $scope.passwordText();
            }
        }, true);

        /**
         * Updates the current connected user information
         * @method saveUserData
         */
        $scope.saveUserData = function(){
            AccountService.saveProfile($scope.item)
                .then(function(data){
                    SessionManagement.currentSession.user = data;
                    $rootScope.currentSession.user = data;

                    if($location.$$search.returnurl){
                        var params = $location.$$search.returnurl.split("?");
                        $location.path(params[0]);
                        if(params[1]) $location.search(params[1]);
                    } else {
                        $location.path("/");
                        $location.search("");
                    }
                });
        }

        /**
         * Updates the password for the current user
         * @method saveNewPassword
         */
        $scope.saveNewPassword = function(){
            AccountService.changeMyPassword($scope.item.id, $scope.password1)
                .then(function() {
                    $scope.password1=null;
                    $scope.password2=null;
                });

        }

        $scope.exit = function(){
            $location.search("");
            $location.path("/");
        };


        $scope.browseIssue = function(i){
            $window.open(BrowserUrlService.issue.edit(i.pkey));
        }


        $scope.gridSubscriptions = IssueService.gridConfig({
            source:IssueService.mySubscriptionsGrid,
            columns: [
                { name: "Key", field:"pkey" },
                { name: "Name", field:"name" },
                { name: "Assignee", field:"assignee" },
                { name: "Reporter", field:"reporter" },
                { name: "Priority", field:"priority" },
                { name: "Status", field:"status" },
                { name: "Due Date", field:"dueDate", filter:"timeAgo" }
            ],
            rowClick:function(r){
                if($scope.clearSubscriptionOnClick){
                    IssueService.subscriptionToggle(r)
                        .then(function(){
                            $scope.gridSubscriptionsParams = {};
                        });
                } else {
                    $window.open(BrowserUrlService.issue.edit(r.pkey));
                }
            }
        });
        $scope.gridSubscriptionsParams = {};
    });
