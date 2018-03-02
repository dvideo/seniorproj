window.onload = function()
{
	document.getElementById('memeid').addEventListener("click", memecontrolax);
}
function memecontrolax(){
	document.getElementById("block1").className = "memeytext";
	document.getElementById("block2").className = "memeytext";
}

