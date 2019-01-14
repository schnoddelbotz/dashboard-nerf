const LOCAL = "Local";
const DASHBOARD = "Dashboard";
const destinations = [LOCAL, DASHBOARD];
var destination = DASHBOARD;

$( document ).ready(function() {

    $("#audio-link").click(function() {
      $("main div.container").hide();
      $("#audio-container").show();
      $("ul.navbar-nav li").removeClass("active");
      $("#audio-li").addClass("active");
    });

    $("#video-link").click(function() {
      $("main div.container").hide();
      $("#video-container").show();
      $("ul.navbar-nav li").removeClass("active");
      $("#video-li").addClass("active");
    });

    $("#speech-link").click(function() {
      $("main div.container").hide();
      $("#speech-container").show();
      $("ul.navbar-nav li").removeClass("active");
      $("#speech-li").addClass("active");
    });

    $("#about-link").click(function() {
      $("main div.container").hide();
      $("#about-container").show();
      $("ul.navbar-nav li").removeClass("active");
      $("#about-li").addClass("active");
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
  } else {
    destination = DASHBOARD;
  }
  updateDestinations();
}
