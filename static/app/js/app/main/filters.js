'use strict';

/* Filters */

angular.module('TrackerApp.filters', [])
    .filter("boolToText", function(){
        return function(value){
            return value ? "YES" : "NO";
        }
    })
    .filter("timeAgo", function(){
        return function(input){
            if(!input){
                return null;
            } else {
                return moment(input).fromNow();
            }
        }
    })
    .filter("percent", function(){
        return function(input){
            if(!input) { return null; }
            return Math.round(parseFloat(input) * 100) + "%";
        }
    })
.filter("dateFormat", function(){
       return function(input){
           if(!input){
               return null;
           } else {
               return moment(input).format("DD-MM-YYYY");
           }
       }
    })
.filter("dateTimeFormat", function(){
        return function(input){
            if(!input){
                return null;
            } else {
                return moment(input).format("DD-MM-YYYY HH:mm");
            }
        }
    })
    .filter("sizeGeneric", function(){
        return function(input){
            var unit = 1024;

            var size=parseFloat(input) || 0;
            var label;
            if(size > unit * unit * unit){
                size = size / (unit * unit * unit);
                label = "Gb";
            } else if (size > unit * unit){
                size = size / (unit * unit);
                label = "Mb";
            } else if (size > unit){
                size = size / unit;
                label = "Kb";
            } else if (size > 0){
                label = "bytes";
            } else {
                return null;
            }
            return size.toFixed(2) + " " + label;
        }
    })
    .filter("trunc", function(){
        return function(text, size){
            if(!text){
                return null;
            } else if(text.length <= size) {
                return text;
            } else {
                return text.substring(0, size) + "...";
            }
        }
    })
    .filter("round", function(){
        return function(value, decimalPositions){
            value = parseFloat(value);
            if(isNaN(value)) return null;
            var factor = 1;
            if(decimalPositions){
                factor = 10 * decimalPositions;
            }
            return Math.round( value * factor) / factor;
        }
    })