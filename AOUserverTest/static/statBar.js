/**
			 * Encapsulates a color in the RGB format for easy
			 * conversion to CSS background color and access
			 * to red, green, and blue values.
			 */
			function RGBColor(r, g, b)
			{
				this.red = r;
				this.green = g;
				this.blue = b;
				this.toCSSBackgroundString = function() {
					return "rgb(" + this.red + "," + this.green + "," + this.blue + ")";
				}
			}
			
			
			// Some global variables.
			var averageRating;
			var numVotes;
			var curRating = 0;
			var ratingIsAdjustable = true;
			
			var startColor = new RGBColor(30, 87, 153);
			var endColors = [
				new RGBColor(68, 150, 182),
				new RGBColor(211, 209, 248),
				new RGBColor(191, 255, 255),
				new RGBColor(191, 255, 223),
				new RGBColor(191, 255, 191),
				new RGBColor(239, 255, 191),
				new RGBColor(255, 239, 191),
				new RGBColor(255, 223, 191),
				new RGBColor(255, 207, 191),
				new RGBColor(255, 191, 191)
			];


			window.onload = function()
			{
				var statContainer = document.getElementById("stat-container");
				var gradientContainer = document.getElementById("gradient-container");
				statContainer.addEventListener("mouseout", hideStatInfo);
				statContainer.addEventListener("mouseover", showStatInfo);
				gradientContainer.addEventListener("mousemove", adjustRating);
				gradientContainer.addEventListener("mousedown", setRating);
				
				// The mouse out listener is triggered even when mousing over a child
				// element, but we only want it to be triggered when mousing out of ths element
				// entirely, so we have to add a bit more code.
				gradientContainer.addEventListener('mouseout', function(event) {
					// Will become true if the mouse was moved over this
					// element or any of its children.
					var isStillOverThisElement = false;
					var thisElement = document.getElementById(this.id);
					var toElement = event.toElement;
					
					// Check to see if the element that was moused over
					// this element or any of its children. If it is,
					// we are still in this element and we don't want
					// to call onMouseUp. This loop traverses through
					// all the parents of the element that was moused
					// over to see if any of them are this element.
					while (toElement != null) {
						if (thisElement == toElement)
							isStillOverThisElement = true;
						toElement = toElement.parentNode;
					}
					
					// Call resetRating if the mouse if not over this element
					// or any of its children.
					if (!isStillOverThisElement)
						resetStat(event);
					}, false);
				
				setupStat();
			}
			
			/**
			 * Makes the stat info bar visible by moving it upward and change the gradient
			 * to a moving bar. 123
			 */
			function showStatInfo()
			{
				document.getElementById("stat-info").style.top='-23px';
			}
			
			/**
			 * Makes the stat info bar invisible by moving it downward and change the moving
			 * bar back into a gradient.
			 */
			function hideStatInfo()
			{
				document.getElementById("stat-info").style.top = '0px'
				$("#gradient-bar").removeClass("rate-line");
				setupStat();
			}
			
			/**
			 * Adjusts the width of the gradient bar and displays the rating
			 * corresponding to the user's mouse position.
			 */
			function adjustRating(event)
			{
				// If the user just voted, he/she cannot adjust the rating without
				// first mousing out of the widget.
				if (ratingIsAdjustable) {
					// Change the gradient bar from a gradient to a line.
					$("#gradient-bar").removeClass("gradient1").addClass("rate-line");
					var gradientBar = document.getElementById("gradient-bar");
					var gradientContainer = document.getElementById("gradient-container");
					
					// Figure out where the mouse is in the stat widget and adjust the size
					// of the gradient and the value of the rating.
					var mouseX = event.pageX - getObjectX(document.getElementById("stat-container"));
					gradientBar.style.width = mouseX + "px";
					var rating = calcRatingFromMouseX(mouseX);
					var endColor = calcEndColorFromRating(rating);
					setBackgroundGradientForElement(document.getElementById("gradient-bar"), startColor, endColor);
					document.getElementById("rating").innerHTML = rating.toFixed(1);
				}
			}
			
			/**
			 * Sets the current rating based on where the user clicks in the widget.
			 * NOTE: Need to include an AJAX request to add the new rating to the database.
			 */
			function setRating(event)
			{
				// If the user just voted, he/she cannot set the rating again
				// without first mousing out of the widget.
				if (ratingIsAdjustable) {
					// Calculate the rating the user just made.
					var mouseX = event.pageX - getObjectX(document.getElementById("stat-container"));
					curRating = calcRatingFromMouseX(mouseX);
					// Update the number of votes and the average rating.
					setupStat();
					// Show the gradient bar as a gradient.
					$("#gradient-bar").removeClass("rate-line");
					var endColor = calcEndColorFromRating(curRating);
					setBackgroundGradientForElement(document.getElementById("gradient-bar"), startColor, endColor);
					ratingIsAdjustable = false;
				}
			}
			
			/**
			 * Sets the gradient bar back to the size of the average,
			 * rating and indicates the user can set the rate again.
			 */
			function resetStat()
			{
				document.getElementById("gradient-bar").style.width = calcMouseXFromRating(averageRating) + "px";
				var endColor = calcEndColorFromRating(averageRating);
				setBackgroundGradientForElement(document.getElementById("gradient-bar"), startColor, endColor);
				ratingIsAdjustable = true;
			}
			
			/**  Sends stat# and calculate average every 10 seconds*/
			function sendStatAvg() {
				setInterval(
					function sendAvg(){ 
						alert("Calculating stat average...");
						resetStat(); 
					}, 3000);
			}
			
			/**
			 * Converts an x coordinate to a rating based on the width of the
			 * gradient container.
			 * @param {double} x - the x coordinate to convert to a rating
			 * @return {double} the rating corresponding to the x coordinate
			 */
			function calcRatingFromMouseX(x)
			{
				// Multiply by 10.1 so that a rating of 10.0 is possible.
				// Add 1 to fix an error I don't fully understand, where if you move the mouse very slowly
				// from the right, the gradient will expand to be the size of the bar, but it won't update
				// the rating to 10.0 (it will still say the average rating).
				return x / parseInt(document.getElementById("gradient-container").offsetWidth + 1) * 10.1;
			}
			
			/**
			 * Converts a rating to an x coordinate based on the width of the
			 * gradient container.
			 * @param {double} r - the rating to convert to an x coordinate
			 * @return {double} the x coordinate corresponding to the rating
			 */
			function calcMouseXFromRating(r)
			{
				// Multiply by 10.1 to make up for the higher rating gotten from the calcRatingFromMouseX.
				// Add one to correspond to the rating from calcRatingFromMouseX.
				return r * parseInt(document.getElementById("gradient-container").offsetWidth + 1) / 10.1;
			}
			
			/**
			 * Sets the average rating and number of votes.
			 * NOTE: Should query database (currently just fills with random numbers).
			 */
			function setupStat()
			{
				/////////////////////////////////////////////////////
				// This part should be replaced by database queries.
				/////////////////////////////////////////////////////
				averageRating = 5;//{{ .avgStat}};
				numVotes = 1; //{{ .numVotes}};
				/////////////////////////////////////////////////////
				// Display the rating and number of votes in the widget.
				document.getElementById("rating").innerHTML = averageRating.toFixed(1);
				document.getElementById("num-votes").innerHTML = numVotes.toLocaleString();
			}
			
			/**
			 * Determines the x position of an HTML element on the page.
			 * Based on code from: http://www.quirksmode.org/js/findpos.html.
			 * @param obj - the object to get the x coordinate from
			 * @return {double} the x coordinate of the object
			 */
			function getObjectX(obj)
			{
				var objX = 0;
				do {
					objX += obj.offsetLeft;
				} while (obj = obj.offsetParent);
				return objX;
			}
			
			/**
			 * Calculates and returns the end color of the gradient
			 * given the rating.
			 * @param {double} rating - the rating to get the color for
			 * @return {RGBColor} the end color of the gradient
			 */
			function calcEndColorFromRating(rating)
			{
				// Get the two colors that the color will be between.
				var lowColor = endColors[Math.floor(rating)];
				var highColor = endColors[Math.floor(rating)];
				
				// This is the position between the two colors.
				// In other words, if the rating is 4.2, we want 0.2
				// of the high color and 0.8 of the low color.
				var pos = rating - Math.floor(rating);
				
				var r = (lowColor.red * (1.0 - pos)) + (highColor.red * pos);
				var g = (lowColor.green * (1.0 - pos)) + (highColor.green * pos);
				var b = (lowColor.blue * (1.0 - pos)) + (highColor.blue * pos);
				
				return new RGBColor(r, g, b);
			}
			
			/**
			 * Sets the background of the given element to a gradient
			 * with the given start and end colors.
			 * @param elem - the element to set the background of
			 * @param {RGBColor} start - the start color
			 * @param {RGBColor} end - the end color
			 */
			function setBackgroundGradientForElement(elem, start, end)
			{
				var startRGB = start.toCSSBackgroundString();
				var endRGB = end.toCSSBackgroundString();
				
				elem.style.background = endRGB; // Old browsers
				elem.style.background = "-moz-linear-gradient(left, " + startRGB + " 0%, " + endRGB + " 100%)"; // FF3.6+
				elem.style.background = "-webkit-gradient(linear, left top, right top, color-stop(0%," + startRGB + "), color-stop(100%," + endRGB + "))"; // Chrome,Safari4+
				elem.style.background = "-webkit-linear-gradient(left, " + startRGB + " 0%," + endRGB + " 100%)"; // Chrome10+,Safari5.1+
				elem.style.background = "-o-linear-gradient(left, " + startRGB + " 0%," + endRGB + " 100%)"; // Opera 11.10+
				elem.style.background = "-ms-linear-gradient(left, " + startRGB + " 0%," + endRGB + " 100%)"; // IE10+
				elem.style.background = "linear-gradient(to right, " + startRGB + " 0%," + endRGB + " 100%)"; // W3C
				elem.style.filter = "progid:DXImageTransform.Microsoft.gradient( startColorstr='" + startRGB + "', endColorstr='" + endRGB + "',GradientType=1 )"; // IE6-9
			}







								var transformProp = Modernizr.prefixed('transform');

					function Carousel3D ( el ) {
						this.element = el;
						this.rotation = 0;
						this.panelCount = 0;
						this.totalPanelCount = this.element.children.length;
						this.theta = 0;

						this.isHorizontal = false;

					}
				
				
					Carousel3D.prototype.modify = function() {

						var panel, angle, i;

						this.panelSize = this.element[ this.isHorizontal ? 'offsetWidth' : 'offsetHeight' ];
						this.rotateFn = this.isHorizontal ? 'rotateY' : 'rotateX';
						this.theta = 360 / this.panelCount;

						// do some trig to figure out how big the carousel
						// is in 3D space
						this.radius = Math.round( ( this.panelSize / 2) / Math.tan( Math.PI / this.panelCount ) );

						for ( i = 0; i < this.panelCount; i++ ) {
							panel = this.element.children[i];
							angle = this.theta * i;
							panel.style.opacity = 1;
							panel.style.backgroundColor = 'hsla(' + angle + ', 100%, 50%, 0.8)';
							// rotate panel, then push it out in 3D space
							panel.style[ transformProp ] = this.rotateFn + '(' + angle + 'deg) translateZ(' + this.radius + 'px)';
						}

						// adjust rotation so panels are always flat
						this.rotation = Math.round( this.rotation / this.theta ) * this.theta;

						this.transform();

					};

					Carousel3D.prototype.transform = function() {
						// push the carousel back in 3D space,
						// and rotate it
						this.element.style[ transformProp ] = 'translateZ(-' + this.radius + 'px) ' + this.rotateFn + '(' + this.rotation + 'deg)';
					};



					var init = function() {


						var carousel = new Carousel3D( document.getElementById('carousel') ),
								panelCountInput = document.getElementById('panel-count'),
								axisButton = document.getElementById('toggle-axis'),
								navButtons = document.querySelectorAll('#navigation button'),

								onNavButtonClick = function( event ){
									var increment = parseInt( event.target.getAttribute('data-increment') );
									carousel.rotation += carousel.theta * increment * -1;
									carousel.transform();
								};

						// populate on startup
						carousel.panelCount = parseInt( 6, 10);
						carousel.modify();


						axisButton.addEventListener( 'click', function(){
							carousel.isHorizontal = !carousel.isHorizontal;
							carousel.modify();
						}, false);

						panelCountInput.addEventListener( 'change', function( event ) {
							carousel.panelCount = event.target.value;
							carousel.modify();
						}, false);

						for (var i=0; i < 2; i++) {
							navButtons[i].addEventListener( 'click', onNavButtonClick, false);
						}

						document.getElementById('toggle-backface-visibility').addEventListener( 'click', function(){
							carousel.element.toggleClassName('panels-backface-invisible');
						}, false);

						setTimeout( function(){
							document.body.addClassName('ready');
						}, 0);

					};
				

					window.addEventListener( 'DOMContentLoaded', init, false);