		// sum of all radio input
		window.sumInputs = function() {
			var inputs = $('.calc:checked'),
		    	result = document.getElementById('total'),
		    	sum = 0;
		for (var i = 0; i < inputs.length; i++) {
		    sum += parseFloat(inputs[i].value);
		}
		$('#total').val(sum);
		
		// in terms of order CO2 tracker weighs the most, followed by water, mileage, and then diet
		// this means that if the user selects the max value for all of the survey question,
		// then they will be led to the CO2 tracker. additionally, the only way for them to get the
		// diet tracker would be if the user didn't select a max value for any of the survey 
		// question besides the first one (the question pertaining to diet)
		var score = sum 
		if (score < 4) {
			alert("Error! Please fill out survey.");			
		} 
		
		else if (score < 100) {
		    alert("Focus on having a greener diet.");
		} 
		
		else if (score < 250) { 
			alert("Focus on reducing your mileage.");
		}
		
		else if (score < 700) { 
			alert("Focus on reducing your water usage.");
		}
		
		else { 
			alert("Focus on reducing CO2 usage.");
		}
			document.getElementById("demo").innerHTML = text;
		}

