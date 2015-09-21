'use strict';



angular.module('TrackerApp.Directives', ['ui.bootstrap'])
    .directive("erpDatePicker", function(){
        return {
            restrict: "E",
            scope:{
                value: '=',
                isOpen: '=',
                isDisabled:'=',
                isRequired:'=',
                minDate:'=',
                maxDate:'='
            },
            controller:function($scope){
                $scope.dateOptions = {
                    'year-format': "'yy'",
                    'starting-day': 1
                };
            },
            templateUrl: '/templates/generic/date_picker.html'
        }
    })/**
 * Clears all behaviors on anchors
 * @method a
 */
    .directive('a', function() {
        return {
            restrict: 'E',
            link: function(scope, elem, attrs) {
                if(attrs.ngClick || attrs.href === '' || attrs.href === '#'){
                    elem.on('click', function(e){
                        e.preventDefault();
                    });
                }
            }
        };
    })
/**
 * Implements SAVE, DELETE and EXIT buttons with icons
 */
    .directive("crudButtons", function(){
        return {
            restrict: 'E',
            templateUrl: '/templates/generic/crud-buttons.html',
            link: function(scope, element, attrs){
            }
        }
    })
    .directive('fileItemDownloadLink', function(PathService){
        return {
            restrict:'E',
            replace:true,
            scope:{
                item:'=',
                forceDownload:'=',
                token:'@'
            },
            link:function(scope, element, attrs){
                if(scope.forceDownload){
                    element.attr('download',scope.item.name);
                }
                scope.url = PathService.file.download;
            },
            template: '<a ng-href="{{url}}/{{item.id}}?token={{token}}" target="_blank">{{item.name}}</a>'
        }
    })
