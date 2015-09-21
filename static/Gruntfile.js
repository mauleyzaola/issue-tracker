var uuid = require('node-uuid'),
    guid = uuid.v4();

module.exports = function(grunt) {

    // 1. All configuration goes here
    grunt.initConfig({
        pkg: grunt.file.readJSON('package.json'),

        clean:{
            main: ['dist/']
        },

        concat: {
            appJs:{
                src:[
                    'app/js/app/app.js',
                    'app/js/app/**/*.js'
                ],
                dest:'dist/scripts/app.js'
            },
            vendorJs:{
                src:[
                    'app/js/lib/jquery/dist/jquery.js',
                    'app/js/lib/angular/angular.js',
                    'app/js/lib/bootstrap/dist/js/bootstrap.js',
                    'app/js/lib/angular-route/angular-route.js',
                    'app/js/lib/ng-file-upload/ng-file-upload.js',
                    'app/js/lib/angular-bootstrap/ui-bootstrap-tpls.js',
                    'app/js/lib/jquery.gritter/js/jquery.gritter.js',
                    'app/js/lib/moment/moment.js'
                ],
                dest:'dist/scripts/vendor.js'
            }
        },

        ngAnnotate:{
            appJs:{
                files:{
                    'dist/scripts/app.js':['dist/scripts/app.js']
                }
            },
            vendorJs:{
                files:{
                    'dist/scripts/vendor.js':['dist/scripts/vendor.js']
                }
            }

        },

        uglify: {
            appJs: {
                src: ['dist/scripts/app.js'],
                dest: 'dist/scripts/app.js'
            },
            vendorJs: {
                src: ['dist/scripts/vendor.js'],
                dest: 'dist/scripts/vendor.js'
            }
        },

        copy:{
            /* css files need to be pointing to the normal versions so we can change at runtime the theme */
            cssGritter:{
                expand:true,
                cwd:'app/js/lib/jquery.gritter/css/',
                src:'jquery.gritter.css',
                dest:'dist/stylesheets/'
            },

            cssFontAwesome:{
                expand:true,
                cwd:'app/js/lib/components-font-awesome/',
                src:'**',
                dest:'dist/stylesheets/font-awesome'
            },

            cssCustom:{
                expand:true,
                cwd:'app/stylesheets/',
                src:'custom.css',
                dest:'dist/stylesheets/'
            },

            index:{
                expand:true,
                cwd:'app/',
                src:'index.html',
                dest:'dist/'
            },
            images:{
                expand:true,
                cwd:'app/images/',
                src:'**',
                dest:'dist/images/'
            },
            templates:{
                expand:true,
                cwd:'app/templates/',
                src:'**',
                dest:'dist/templates/'
            },
            faFonts:{
                expand:true,
                cwd:'app/stylesheets/themes/font-awesome/fonts/',
                src:'**',
                dest:'dist/fonts/'
            },
            bootstrapFonts:{
                expand:true,
                cwd:'app/js/lib/bootstrap/fonts/',
                src:'**',
                dest:'dist/fonts/'
            },
            imagesGritter:{
                expand:true,
                cwd:'app/js/lib/jquery.gritter/images/',
                src:'**',
                dest:'dist/images/'
            },
            underscore:{
                expand:true,
                cwd:'app/js/lib/underscore/',
                src:'underscore-min.js',
                dest:'dist/scripts/'
            },
            angulari18n:{
                expand:true,
                cwd:'app/js/lib/angular-i18n/',
                src:'angular-locale_es-mx.js',
                dest:'dist/scripts/'
            },
            angulariSanitize:{
                expand:true,
                cwd:'app/js/lib/angular-sanitize/',
                src:'angular-sanitize.min.js',
                dest:'dist/scripts/'
            },
            momentLocale:{
                expand:true,
                cwd:'app/js/lib/moment/locale/',
                src:'es.js',
                dest:'dist/scripts/'
            }
        },

        processhtml:{
            build:{
                files:{
                    'dist/index.html':['dist/index.html']
                }
            }
        },

        'string-replace':{
            inline:{
                files:{
                    'dist/index.html': 'dist/index.html'
                },
                options:{
                    replacements:[
                        {
                            pattern: '<!--start DEV css imports-->',
                            replacement: '<!--start DEV css imports'
                        },
                        {
                            pattern: '<!--end DEV css imports-->',
                            replacement: 'end DEV css imports-->'
                        },
                        {
                            pattern: '<!--start PROD css imports',
                            replacement: '<!--start PROD css imports-->'
                        },
                        {
                            pattern: 'end PROD css imports-->',
                            replacement: '<!--end PROD css imports-->'
                        },
                        {
                            pattern: 'app.js',
                            replacement: 'app.js?v=' + guid
                        },
                        {
                            pattern: 'vendor.js',
                            replacement: 'vendor.js?v=' + guid
                        },
                        {
                            pattern: 'main.css',
                            replacement: 'main.css?v=' + guid
                        },
                        {
                            pattern: '/js/lib/underscore/underscore.js',
                            replacement: '/scripts/underscore-min.js'
                        },
                        {
                            pattern: '/js/lib/angular-i18n/angular-locale_es-mx.js',
                            replacement: '/scripts/angular-locale_es-mx.js'
                        },
                        {
                            pattern: '/js/lib/angular-sanitize/angular-sanitize.js',
                            replacement: '/scripts/angular-sanitize.min.js'
                        },
                        {
                            pattern: '/js/lib/moment/locale/es.js',
                            replacement: '/scripts/es.js'
                        },
                        {
                            pattern: 'Google Analytics -->',
                            replacement: '<!-- Google Analytics -->'
                        },
                        {
                            pattern: '<!-- Google Analytics',
                            replacement: '<!-- Google Analytics -->'
                        }
                    ]
                }
            }
        }
    });

    // 3. Where we tell Grunt we plan to use this plug-in.
    grunt.loadNpmTasks('grunt-contrib-concat');
    grunt.loadNpmTasks('grunt-contrib-uglify');
    grunt.loadNpmTasks('grunt-ng-annotate');
    grunt.loadNpmTasks('grunt-contrib-copy');
    grunt.loadNpmTasks('grunt-contrib-clean');
    grunt.loadNpmTasks('grunt-processhtml');
    grunt.loadNpmTasks('grunt-string-replace');


    // 4. Where we tell Grunt what to do when we type "grunt" into the terminal.
    grunt.registerTask('default', [
        'clean',
        'concat',
        'ngAnnotate',
        'uglify',
        'copy',
        'processhtml',
        'string-replace'
    ]);

};
