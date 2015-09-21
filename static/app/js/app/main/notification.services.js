'use strict';

angular.module("TrackerApp.Notification.services", [])
    .factory("NotificationTypes", function(){
        return {

            objectType: {
                attachment:"attachment",
                group:"group",
                permissionScheme:"permissionScheme",
                priority:"priority",
                profile:"profile",
                project:"project",
                issue:"issue",
                issueAttachment:"issueAttachment",
                issueComment:"issueComment",
                role:"role",
                session:"session",
                user:"user",
                status:"status",
                workflow:"workflow",
                workflowStep:"workflowStep"
            },

            operation:{
                add:"insert",
                update:"update",
                delete:"delete"
            },

            iconType:{
                ok:"ok",
                comment:"comment",
                error:"error",
                info:"info",
                chat:"chat",
                warning:"warning"
            }
        }
    })

    .factory("NotificationService", function(BrowserUrlService, Notifier, $window, NotificationTypes){
        var createNotificationText= function(data){
                if(!data || !data.objectType || !data.operation){
                    $window.console.log("Object type and operation are necessary to create a notification");
                    return;
                }

                var title, url, text, textObjectName, textObjectType;

                switch (data.operation){
                    case NotificationTypes.operation.add:
                        title="Data Added";
                        break;
                    case NotificationTypes.operation.delete:
                        title="Data Removed";
                        break;
                    case NotificationTypes.operation.update:
                        title="Data Updated";
                        break;
                    default :
                        throw ("Unknown Operation");
                        break;
                }


                switch (data.objectType){

                    case NotificationTypes.objectType.group:
                        textObjectType="Group";
                        if(data.item){
                            if(data.item.id && data.operation != NotificationTypes.operation.delete){
                                url=BrowserUrlService.group.edit(data.item.id);
                            }
                            textObjectName=data.item.name;
                        }

                        break;

                    case NotificationTypes.objectType.role:
                        textObjectType="Role";
                        if(data.item){
                            if(data.item.id && data.operation != NotificationTypes.operation.delete){
                                url=BrowserUrlService.role.edit(data.item.id);
                            }
                            textObjectName=data.item.name;
                        }

                        break;


                    case NotificationTypes.objectType.user:
                        textObjectType="User";
                        if(data.item){
                            if(data.item.id && data.operation != NotificationTypes.operation.delete){
                                url=BrowserUrlService.user.edit(data.item.id);
                            }
                            textObjectName=data.item.name + " " + data.item.lastName;
                        }

                        break;

                    case NotificationTypes.objectType.session:
                        textObjectType="Session";
                        break;


                    case NotificationTypes.objectType.issue:
                        textObjectType="Issue";
                        if(data.item){
                            if(data.item.id && data.operation != NotificationTypes.operation.delete){
                                url=BrowserUrlService.issue.edit(data.item.pkey);
                            }
                            textObjectName="(" + data.item.pkey + ") " + data.item.name;
                        }

                        break;


                    case NotificationTypes.objectType.issueAttachment:
                    case NotificationTypes.objectType.issueComment:
                        if(data.objectType === NotificationTypes.objectType.issueAttachment){
                            textObjectType="Attachment";
                        } else if(data.objectType === NotificationTypes.objectType.issueComment){
                            textObjectType="Comment";
                        }
                        if(data.item && data.operation){
                            url = BrowserUrlService.issue.edit(data.item.issue.pkey);
                            textObjectName = "(" + data.item.issue.pkey + ") " + data.item.issue.name + "<br/>";
                            switch (data.objectType){
                                case NotificationTypes.objectType.issueComment:
                                    textObjectName += data.item.body;
                                    break;
                                case NotificationTypes.objectType.issueAttachment:
                                    textObjectName += data.item.fileItem.name;
                                    break;
                            }
                        }
                        break;

                    case NotificationTypes.objectType.permissionScheme:
                        textObjectType="Permission Scheme";
                        if(data.item){
                            if(data.item.id && data.operation != NotificationTypes.operation.delete){
                                url=BrowserUrlService.permissionScheme.edit(data.item.id);
                            }
                            textObjectName=data.item.name;
                        }

                        break;


                    case NotificationTypes.objectType.priority:
                        textObjectType="Priority";
                        if(data.item){
                            if(data.item.id && data.operation != NotificationTypes.operation.delete){
                                url=BrowserUrlService.priority.edit(data.item.id);
                            }
                            textObjectName=data.item.name;
                        }

                        break;

                    case NotificationTypes.objectType.profile:
                        textObjectType="User Profile";
                        if(data.item && data.item.name){
                            textObjectName=data.item.name + ' ' + data.item.lastName;
                        }
                        break;

                    case NotificationTypes.objectType.project:
                        textObjectType="Project";
                        if(data.item){
                            if(data.item.id && data.operation != NotificationTypes.operation.delete){
                                url=BrowserUrlService.project.edit(data.item.id);
                            }
                            textObjectName=data.item.name;
                        }

                        break;

                    case NotificationTypes.objectType.status:
                        textObjectType="Status";
                        break;

                    case NotificationTypes.objectType.workflow:
                        textObjectType="Workflow";
                        if(data.item){
                            if(data.item.id && data.operation != NotificationTypes.operation.delete){
                                url=BrowserUrlService.workflow.edit(data.item.id);
                            }
                            textObjectName=data.item.name;
                        }

                        break;
                    case NotificationTypes.objectType.workflowStep:
                        textObjectType="Workflow Step";
                        if (data.item && data.item.name) {
                            textObjectName = data.item.name;
                        }
                        break;
                }


                if(!textObjectType){ return; }

                if(url && textObjectName){
                    text=textObjectType + ": <a href='" + url + "'>" + textObjectName + "</a>";
                } else if(textObjectName){
                    text=textObjectType + ":" + textObjectName;
                } else{
                    text=textObjectType;
                }

                return {
                    title:title,
                    text:text
                }
            }

        return {
            notify:function(data){
                var notificationText = createNotificationText(data);
                if(!notificationText){
                    return;
                }
                Notifier.push(notificationText);
            },

            notificationText:function(data){
                return createNotificationText(data);
            },

            resolveNotificationIconUrl:function(data){
                var icon;
                try{
                    icon = NotificationTypes.iconType[data];
                    if(!icon){
                        throw("");
                    }
                }catch(e){
                    icon = "/images/" + NotificationTypes.iconType.info + ".png";
                }
                return "/images/" + icon + ".png";
            }
        }
    });
