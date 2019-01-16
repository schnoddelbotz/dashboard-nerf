const LOCAL = "Local";
const DASHBOARD = "Dashboard";
const destinations = [LOCAL, DASHBOARD];
var destination = DASHBOARD;

$( document ).ready(function() {

    $("#videofilter").on("keyup", function(){
      var value = $(this).val().toLowerCase();
      $("#video-container ul.videos li").filter(function() {
        $(this).toggle($(this).text().toLowerCase().indexOf(value) > -1)
      });
    });

    const handleNav = function(selectedMedia) {
      $("main div.container").hide();
      $("#" + selectedMedia + "-container").show();
      $("ul.navbar-nav li").removeClass("active");
      $("#" + selectedMedia + "-li").addClass("active");
    }

    $("#audio-link").click(function() {
      handleNav("audio");
    });

    $("#video-link").click(function() {
      handleNav("video");
    });

    $("#speech-link").click(function() {
      handleNav("speech");
    });

    $("#about-link").click(function() {
      handleNav("about");
    });

    $("#stop-link").click(function() {
      $.get("stop", function( data ) {});
    });

    $("ul.songs li").click(function(e) {
      play('audio', e.target.innerHTML);
    });

    $("ul.videos li").click(function(e) {
      play('video', e.target.innerHTML);
    });

    $("#speak-button").click(function(e) {
      speak($("#speak-text").val());
    });
    $("#speak-text").on('keyup', function (e) {
      if (e.keyCode == 13) {
          speak($("#speak-text").val());
      }
    });

    updateDestinations();

});

function play(mediaType, filename) {
  if (destination == LOCAL) {
    var player = mediaType + 'Player';
    $("#"+player).show();
    var el = document.getElementById(player),
    elClone = el.cloneNode(true);
    el.parentNode.replaceChild(elClone, el);
    var media = $("#"+player);
    $("#"+mediaType).attr("src", "media/"+filename);
    media[0].load();
    media[0].oncanplaythrough = media[0].play();
    media[0].onended = function(){ $("#"+player).hide(); }
  } else {
    $.get( "play/"+mediaType+"/"+filename, function( data ) {});
  }
}

function speak(text) {
  if (destination == LOCAL) {
    window.speechSynthesis.speak(new SpeechSynthesisUtterance(text));
  } else {
    $.get( "speech/?text="+text, function( data ) {});
  }
}

function updateDestinations() {
  $("#destinations").html("");
  $.each(destinations, function( key, val ) {
    var cssClass = "dropdown-item";
    var oc = 'onclick="setDestination(\''+val+'\'); return false;"'
    if (destination == val) {
      cssClass = cssClass + " active";
    }
    $("#destinations").append('<a href="#" '+oc+' class="'+cssClass+'">'+val+'</a>');
  });
}

function setDestination(dest) {
  if (dest==LOCAL) {
    destination = LOCAL;
    $("#stop-li").hide();
  } else {
    destination = DASHBOARD;
    $("#stop-li").show();
  }
  updateDestinations();
}
