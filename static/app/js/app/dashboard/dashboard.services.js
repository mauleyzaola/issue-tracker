'use strict';

angular.module("TrackerApp.Dashboard.services", [])
    .factory("DashboardService", function($http, PathService, RunApiService){
        return {
            finder:{
                search:function(data){
                    return $http.get(RunApiService.generateUrl(PathService.services.search,data))
                        .then(function(response){
                            var data = response.data;
                            if(!data.rows){
                                data.rows = [];
                            }
                            for(var i=0; i < data.rows.length; i++){
                                data.rows[i] = JSON.parse(data.rows[i]);
                            }
                            return data;
                        });
                }
            },

            issue:{
                groupAll:function(data){
                    return $http.get(PathService.issue.groupAll(data))
                        .then(function(response){
                            return response.data;
                        });
                },

                groupByDataType:function(data){
                    return $http.get(RunApiService.generateUrl(PathService.issue.groupByDataType, data))
                        .then(function(response){
                            return response.data;
                        });
                }
            }
        }
    });