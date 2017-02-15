// Open specified tab
function openTab(evt, tabName) {
    var i, tabcontent, tablinks, info;
	info = document.getElementById("Info");
	info.style.display = "none";
    tabcontent = document.getElementsByClassName("tabcontent");
    for (i = 0; i < tabcontent.length; i++) {
        tabcontent[i].style.display = "none";
    }
    tablinks = document.getElementsByClassName("tablinks");
    for (i = 0; i < tablinks.length; i++) {
        tablinks[i].className = tablinks[i].className.replace(" active", "");
        tablinks[i].className = tablinks[i].className.replace(" inactive", "");
        tablinks[i].className += " inactive";
    }
    document.getElementById(tabName).style.display = "block";
    evt.currentTarget.className = evt.currentTarget.className.replace(" inactive", "");
    evt.currentTarget.className += " active";
}
