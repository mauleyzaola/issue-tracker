'use strict';

angular.module('TrackerApp.Grid.Directives', [])
    .directive('erpGrid', function($compile) {

        return {
            restrict: 'E',
            replace:true,
            scope:{

                config:'=config',

                params:'=params',

                rows:'@'
            },

            link:function(scope, element, attributes){

                var parseValue = function(lastValue, field){
                    var objectPath = field.split(".");
                    if(!objectPath) { return null; }
                    else if (objectPath.length == 1) {
                        return lastValue[objectPath[0]];
                    }
                    else {
                        var nextValue = lastValue[objectPath[0]];
                        if(nextValue == null){
                            return nextValue;
                        }
                        objectPath.splice(0,1);
                        return parseValue(nextValue, objectPath.join("."));
                    }
                };

                scope.getValue = function(row, col){
                    return parseValue(row, col);
                };


                scope.$watch('config.columns', function(newVal){
                    if(!newVal){ return; }
                    var col,
                        html = '<div>';
                    html += '<div class="row">';
                    html += '<div class="col-lg-12">';
                    html += '<table class="' + scope.config.customCss + '">';
                    html += '<thead>';
                    html += '<tr>';
                    html += '<th ng-class="col.css" ng-repeat="col in config.columns" ng-click="sortColumn(col)" >{{col.name}} ';
                    html += '<i class="fa" ng-class="col.class"></i>';
                    html += '</th>';
                    html += '</tr>';
                    html += '</thead>';
                    html += '<tbody>';
                    html += '<tr ng-repeat="r in rows" ng-click="config.rowClick(r)">';
                    for (var i=0;i<newVal.length; i++){
                        var field="";
                        col = newVal[i];
                        html += '<td>';
                        if(col.filter){
                            field += '{{getValue(r,"' + col.field + '") | ' + col.filter + '}}';
                        } else {
                            field += '{{getValue(r,"' + col.field + '")}}';
                        }
                        if(col.template){
                            html += col.template.replace('@', field);
                        } else {
                            html += field;
                        }
                        html += '</td>';
                    }
                    html += '</tr>';
                    html += '</tbody>';
                    html += '</table>';
                    html += '</div>';
                    html += '</div>';

                    html += '<div class="row">';
                    html += '<div class="col-lg-3">';
                    html += '<div class="pull-right input-group merged">';
                    html += '<span class="input-group-addon"><i class="fa fa-search"></i> </span>';
                    html += '<input class="form-control" type="text" ng-model="config.grid.query" />';
                    html += '</div>';
                    html += '</div>';
                    html += '<div class="col-lg-3">';
                    html += '<span>Page</span>';
                    html += '<select ng-model="config.grid.pageSize" ng-options="option as option for option in config.grid.pageSizes"></select>';
                    html += '</div>';
                    html += '<div class="col-lg-6">';
                    html += '<div class="pull-right">';
                    html += '<ul class="pagination pagination-sm">';
                    html += '<li ng-repeat="page in config.grid.pages" ng-class="{active: page.active, disabled: page.disabled}"><a href="#" ng-click="changePage(page)">{{page.title}}</a></li>';
                    html += '</ul>';
                    html += '</div>';
                    html += '</div>';
                    html += '</div>';

                    html += '</div>';
                    element.html($compile(html)(scope));
                })
            },

            controller:function($scope){

                if($scope.config === undefined){ $scope.config = {}; }

                if($scope.config.grid === undefined){ $scope.config.grid = {}; }
                if($scope.config.grid.sortDirection === undefined){ $scope.config.grid.sortDirection = 'asc'; }
                if($scope.config.grid.pageSizes === undefined){ $scope.config.grid.pageSizes = [5,10,20,50]; }
                if($scope.config.grid.pageSize === undefined){ $scope.config.grid.pageSize = 10; }
                if($scope.config.grid.pageNumber === undefined){ $scope.config.grid.pageNumber = 1; }
                if($scope.config.grid.maxDisplayPages === undefined){ $scope.config.grid.maxDisplayPages = 5; }

                var sortDirection = $scope.config.grid.sortDirection == 'asc';

                /**
                 * Validate changes on properties
                 * @param newVal
                 * @param oldVal
                 */
                var checkIsDirty = function(newVal, oldVal){
                    if(newVal != oldVal) {
                        if($scope.params){
                            $scope.loadData();
                        }
                    }
                }

                var setPages = function(){
                    var pages = [];
                    if(!$scope.config.grid.totalCount || $scope.config.grid.totalCount == 0){
                        return pages;
                    }

                    //first and previous
                    var pageNumber=$scope.config.grid.pageNumber, pageCount = $scope.config.grid.pageCount, maxDisplayPages = $scope.config.grid.maxDisplayPages;
                    pages.push({ title: "First", index: 1, disabled: pageNumber <= 1 });

                    if(pageNumber > 1) { pages.push({ title:"Prev", index: pageNumber-1 }); }


                    //previous set
                    if(pageNumber > maxDisplayPages){
                        pages.push({ title: "...", index: pageNumber - 5});
                    }

                    //numbered pages
                    for(var i=0; i < maxDisplayPages; i++){
                        var index = i+pageNumber;
                        if(index <= pageCount){ pages.push({ title: index, index: index, active: index == pageNumber}); }
                    }

                    //next set
                    if(pageNumber < pageCount - maxDisplayPages){
                        pages.push({ title: "...", index: pageNumber + 5});
                    }

                    //next and last
                    if(pageNumber < pageCount) { pages.push({ title:"Next", index: pageNumber+1 }); }
                    pages.push({ title: "Last", index: pageCount, disabled: pageNumber >= pageCount });

                    return pages;
                };

                /**
                 * Combines internal and custom paramters to send to load service
                 */
                var getSourceParams = function(){
                    var firstColumn='';
                    if($scope.config.columns.length != 0){ firstColumn = $scope.config.columns[0].field; }
                    var pars = {};
                    angular.copy($scope.params, pars);
                    pars.sortField = $scope.config.grid.sortField || firstColumn;
                    pars.sortDirection = $scope.config.grid.sortDirection;
                    pars.pageSize = $scope.config.grid.pageSize;
                    pars.pageNumber = $scope.config.grid.pageNumber;
                    pars.query = $scope.config.grid.query || '';
                    return pars;
                };

                $scope.loadData = function(){
                    $scope.config.source(getSourceParams())
                        .then(function(data){
                            var rows = data.rows;
                            $scope.rows = rows;
                            $scope.config.grid.totalCount = data.totalCount;
                            $scope.config.grid.pageNumber = data.pageNumber;

                            if(data.totalCount == 0) {
                                $scope.config.grid.pageCount = 0;
                            } else {
                                var pageCount = parseInt(data.totalCount / $scope.config.grid.pageSize),
                                    remainder = data.totalCount % $scope.config.grid.pageSize;
                                pageCount += Math.sign(remainder);
                                $scope.config.grid.pageCount = pageCount;
                            }
                            $scope.config.grid.pages = setPages();
                        });
                }

                $scope.$watch(function(){ return $scope.params; }, function(newVal){
                    if(newVal !== undefined){
                        if($scope.config.columns && $scope.config.columns.length != 0){
                            $scope.loadData();
                        }
                    }
                });

                $scope.sortColumn = function(col){
                    if(col && $scope.params){
                        sortDirection = !sortDirection;
                        $scope.config.grid.sortDirection = sortDirection ? 'asc' : 'desc';
                        $scope.config.grid.sortField = col.field;

                        for(var i=0;i<$scope.config.columns.length; i++){
                            var c = $scope.config.columns[i];
                            if(c.field === col.field){
                                c.class = !sortDirection ?  'fa-sort-desc' : 'fa-sort-asc';
                            } else {
                                c.class = "";
                            }
                        }
                        $scope.loadData();
                    }
                };


                $scope.changePage = function(page){
                    $scope.config.grid.pageNumber = page.index;
                };

                $scope.$watch('config.grid.pageNumber', checkIsDirty);
                $scope.$watch('config.grid.pageSize', checkIsDirty);
                $scope.$watch('config.grid.query', checkIsDirty);
            }
        };
    });