'use strict';

angular.module('TrackerApp.Issue.Directives', [])
    .directive("issueList", function(){
        return {
            restrict: "E",
            replace: true,
            scope: {
                items: '=items',
                resolveUrl: "=resolveUrl"
            },
            templateUrl: "/templates/issue/issue/issue_list.html"
        }
    })
    .directive("issueItem", function(){
        return {
            restrict: "E",
            replace: true,
            scope: {
                item: '=item',
                resolveUrl: "=resolveUrl"
            },
            templateUrl: "/templates/issue/issue/issue_item.html"
        }
    })
    .directive("issueCondensedItem", function(){
        return {
            restrict: "E",
            replace: true,
            scope: {
                item: '=item',
                resolveUrl: "=resolveUrl"
            },
            templateUrl: "/templates/issue/issue/issue_condensed_item.html"
        }
    })