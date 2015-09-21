'use strict';

angular.module("TrackerApp.Account.services", [])
    .factory("AccountService", function($http, $q, $window, Notifier, PathService, NotificationService, NotificationTypes){
        return {
            login: function(data){
                return $http.post(PathService.account.login, data)
                    .then(function(response){
                        return response.data;
                    });
            },

            logout: function(){
                return $http.post(PathService.account.logout)
                    .then(function(response){
                        return response.data;
                    });
            },

            saveProfile: function(data){
                return $http.put(PathService.account.profile, data)
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.profile,
                            operation: NotificationTypes.operation.update,
                            item: response.data
                        });


                        return response.data;
                    });
            },

            loadProfile: function(){
                return $http.get(PathService.account.profile)
                    .then(function(response){
                        return response.data;
                    });
            },

            removeToken: function(token){
                return $http.delete(PathService.account.removeToken(token))
                    .then(function(response){
                        NotificationService.notify({
                            objectType: NotificationTypes.objectType.session,
                            operation: NotificationTypes.operation.delete,
                            item: response.data
                        });

                        return response.data;
                    });
            },

            getTokens: function(){
                return $http.get(PathService.account.getTokens)
                    .then(function(response){
                        return response.data;
                    });
            },

            passwordRecoverToken: function(email){
                return $http.post(PathService.account.passwordRecoverToken, { email: email})
                    .then(function(response){
                        Notifier.alert({
                            title: "Request Sent",
                            text: "You should recieve soon an email with instructions on how to reset your password"
                        });
                        return response.data;
                    });
            },

            verifyToken: function(token){
                return $http.post(PathService.account.verifyToken, {tokenEmail:token})
                    .then(function(response){
                        return response.data;
                    });
            },

            changeMyPassword: function(id, newPassword){
                return $http.post(PathService.account.changeMyPassword, { id: id, password:newPassword})
                    .then(function(response){
                        Notifier.push({title:"Password change", text:"Password has been changed successfully"});
                        return response.data;
                    });
            },


            changePassword: function(id, newPassword){
                return $http.post(PathService.account.changePassword, { id: id, password:newPassword})
                    .then(function(response){
                        Notifier.push({title:"Password change", text:"Password has been changed successfully"});
                        return response.data;
                    });
            },

            delaySession:function(){
                return $http.post(PathService.account.delaySession)
                    .then(function(response){
                        return response.data;
                    });
            }
        }
    });
