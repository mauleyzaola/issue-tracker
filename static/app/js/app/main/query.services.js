'use strict';

angular.module("TrackerApp.QueryString.services", [])

    .factory("QueryStringNames", function(){
        return {
            assignee: "assignee",
            group:"group",
            formula:"formula",
            id: "id",
            isCancelled: "isCancelled",
            parent: "parent",
            pkey: "pkey",
            priority: "priority",
            project: "project",
            reporter: "reporter",
            resolved: "resolved",
            status: "status",
            target: "target",
            token: "token",
            workflow: "workflow"
        }
    })
