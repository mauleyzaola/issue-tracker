'use strict';





angular.module('TrackerApp.services', [])
    .factory('dashboardFn', function(){
        return {
            cssColor:function(percentValue){
                if(percentValue > 0.8){
                    return 'success';
                } else if(percentValue > 0.6){
                    return 'default';
                } else if(percentValue > 0.4){
                    return 'info';
                } else if(percentValue > 0.2){
                    return 'warning';
                } else {
                    return 'danger';
                }
            }
        }
    })
    .factory("utils", function($window){
        return {
            guid: function(){
                function S4() {
                    return (((1+Math.random())*0x10000)|0).toString(16).substring(1);
                }
                // then to call it, plus stitch in '4' in the third group
                return (S4() + S4() + "-" + S4() + "-4" + S4().substr(0,3) + "-" + S4() + "-" + S4() + S4() + S4()).toLowerCase();
            },

            replaceAll: function(text, search, newText){
                var re=new RegExp(search, 'g');
                return text.replace(re, newText);
            },

            confirm: function(text){
                return $window.confirm(text || "Sure to remove this data?");
            },

            isNumber: function(value){
                return !isNaN(parseFloat(value)) && isFinite(value);
            },

            findIndexById: function(id, elements){
                if(!id || !elements || elements.length==0){ return -1; }
                for(var i=0; i<elements.length; i++){
                    if(elements[i].id == id){
                        return i;
                    }
                }
                return -1;
            },

            queryStringFromObject:function(params){
                var qs = function(obj, prefix){
                    var str = [];
                    for (var p in obj) {
                        var k = prefix ? prefix + "[" + p + "]" : p,
                            v = obj[p];
                        str.push(angular.isObject(v) ? qs(v, k) : (k) + "=" + encodeURIComponent(v));
                    }
                    return str.join("&");
                }
                return qs(params);
            }
        }
    })
    .factory("SessionManagement", function($window){

        var _path = "currentSession";
        var _value;


        var loadCurrentSession = function(){
            try {
                if(_value){
                    return _value;
                }
                var value = $window.localStorage[_path];
                if(!value){
                    return null;
                }
                value = $window.atob(value);
                if(value){
                    return JSON.parse(value);
                }
            }catch (e){
                $window.console.log(e);
                _value=null;
                return null;
            }
        }

        var saveCurrentSession = function(value){
            if(value){
                _value=value;
                $window.localStorage[_path] = $window.btoa(JSON.stringify(value));
                return _value;
            }
        }

        return {
            currentSession: function(value){
                if(value){
                    return saveCurrentSession(value);
                } else {
                    return loadCurrentSession();
                }
            },

            clearCurrentSession: function(){
                _value=null;
                $window.localStorage.removeItem(_path);
            }
        }
    })
    .factory("DefaultStyles", function(){
        return {
            css:{
                defaultTableHoverCss:'table table-hover table-bordered table-striped',
                defaultTableCss:'table table-bordered table-striped'
            },
            templates:{
                rightCell: '<span class="pull-right">@</span> '
            }
        }
    })
    .factory("Notifier", function(){
        return {
            push: function(options){
                $.gritter.options.position = "bottom-right";
                $.gritter.add({
                    title: options.title,
                    text: options.text,
                    image: options.image || "/images/ok.png"
                });
            },

            warning: function(options){
              $.gritter.options.position = "bottom-right";
              $.gritter.add({
                title: options.title || "Warning",
                text: options.text || "Warning",
                sticky: false,
                image: options.image || "/images/warning.png"
              });
            },
            info: function(options){
                $.gritter.options.position = "bottom-right";
                $.gritter.add({
                    title: options.title || "Info",
                    text: options.text || "Info",
                    sticky: false,
                    image: options.image || "/images/info.png"
                });
            },
            error: function(options){
                $.gritter.options.position = "bottom-right";
                $.gritter.add({
                    title: options.title || "Error",
                    text: options.text || "Error",
                    sticky: false,
                    image: options.image || "/images/error.png"
                });
            },


            alert: function(options){
                if($.pnotify){
                    var delay = options.delay || 2000;
                    if(options.isError){
                        delay = 4000;
                    }
                    var pWidth = options.width || "50%";

                    var pinType;
                    if(options.isError){
                        pinType="error";
                        if(!options.title){
                            options.title = "Error";
                        }
                    }
                    else{
                        pinType= options.isSuccess ? "success" : "info";
                        if(!options.title){
                            options.title= options.isSuccess ? "Ok" : "Information";
                        }
                    }
                    $.pnotify({
                        title: options.title,
                        text: options.text,
                        type: options.pinType,
                        width: options.Width,
                        nonblock: false,
                        history: false,
                        animate_speed: "slow",
                        addclass: "stack-bar-top",
                        delay: options.delay
                    });
                }
                else{
                    alert(options.text);
                }
            }
        }
    });
