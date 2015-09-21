'use strict';

angular.module('TrackerApp.Account.Directives', [])
    .directive("notificationItems", function(){
        return {
            restrict: "E",
            replace: true,
            scope: {
                items: '=items',
                resolveUrl: "=resolveUrl",
                removeItem: "=removeItem"
            },
            templateUrl: "/templates/account/notification_items.html"
        }
    })