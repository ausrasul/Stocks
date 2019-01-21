module.exports = function(grunt) {

// Project configuration.
	grunt.initConfig({
		pkg: grunt.file.readJSON('package.json'),
		uglify: {
			options: {
				banner: '/*! Base <%= grunt.template.today("yyyy-mm-dd") %> */\n'
			},
			dist: {
				files: {
					'static/js/base.min.js': ['<%= concat.dist.dest %>']
				}
			}
		},
		concat: {
			options: {
				// define a string to put between each file in the concatenated output
				separator: ';'
			},
			dist: {
				// the files to concatenate
				src: [
					'js/src/primal.js',
					'js/src/tools.js',
					'js/src/muiInit.js',
					'js/src/infoBar.js',
				],
				// the location of the resulting JS file
				dest: 'static/js/base.js'
			}
		},
		jshint: {
			files: ['Gruntfile.js', 'js/src/**/*.js', 'src/js/tests/spec/*.js'],
			options: {
				// options here to override JSHint defaults
				globals: {
					console: true,
					document: true
				}
			}
		},
		sass: {                              // Task 
			dist: {                            // Target 
	  			options: {                       // Target options 
					style: 'compressed',
					sourcemap: 'inline',
					loadPath: ['sass/']
	 			},
	  			files: {                         // Dictionary of files 
					'static/css/base.css': 'sass/base.scss'
	  			}
			}
  		},
		watch: {
			scripts: {
				files: ['<%= jshint.files %>'],
				tasks: [
					'jshint',
					'concat:dist'],
				options: {
						interrupt: true,
						livereload: {
							host: 'devserver.localhost',
							port: 9000,
							key: grunt.file.read('certs/test.key'),
							cert: grunt.file.read('certs/test.crt')
        					// you can pass in any other options you'd like to the https server, as listed here: http://nodejs.org/api/tls.html#tls_tls_createserver_options_secureconnectionlistener
      					}
				}
			},
			sass: {
				files: ['sass/*.scss'],
				tasks: ['sass'],
				options: {
					interrupt: true,
					livereload: true
				}	
			},
			twig: {
				files: ['views/*.twig'],
				options: {
					livereload: true
				}
				
			}
		}
	
	
	});

	// Load the plugin that provides the "uglify" task.
	grunt.loadNpmTasks('grunt-contrib-uglify');
	grunt.loadNpmTasks('grunt-contrib-watch');
	grunt.loadNpmTasks('grunt-contrib-concat');
	grunt.loadNpmTasks('grunt-contrib-jshint');
	grunt.loadNpmTasks('grunt-contrib-sass');
	// Default task(s).
	grunt.registerTask('default', ['jshint', 'concat', 'uglify', 'sass']);
	//grunt.registerTask('js', ['jshint', 'concat']);
	//grunt.registerTask('css', ['jshint', 'less']);

};
