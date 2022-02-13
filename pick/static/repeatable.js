const requiredParticipants = 3
// To add javascript functionality - if javascript is enabled
document.addEventListener("DOMContentLoaded", function () {

  startupParticipants(requiredParticipants)

  const form = document.querySelector('#secretSantaForm');
  form.addEventListener('submit', handleSubmit);
});

function startupParticipants(n) {
  let people = document.getElementById('people');
  while (people.childElementCount < n) {
    addPerson()
  }
}

// add a new person to the end of the form when "Add Person" is clicked
function addPerson() {
  let people = document.getElementById('people');
  let person = people.children[0];

  // clonedNode doesn't copy event listeners
  var clone = person.cloneNode(true);

  // clear out the textfields for the new person. Don't copy the person info
  inputs = clone.getElementsByTagName('input');
  for (i = 0; i < inputs.length; i++) {
    if (inputs[i].type == "text" || inputs[i].type == "email") {
      inputs[i].value = '';
    }
  }

  if (people.childElementCount > requiredParticipants) {
    // remove button already exists. Copy the event listener since clone
    // doesn't copy that
    btnCount = clone.getElementsByClassName("remove-person");
    btnCount[0].onclick = removePerson
  }

  people.appendChild(clone);



  if (people.childElementCount == requiredParticipants + 1) {
    // has increased above required for the first time, add remove button to all
    addRemovePerson()
  }
}


// addRemovePerson adds the Remove Person button to each person
function addRemovePerson() {
  let people = document.getElementById('people');
  personCount = people.childElementCount

  for (var i = 0; i < personCount; i++) {
    person = people.children[i];
    btnCount = person.getElementsByClassName("remove-person");

    let btn = document.createElement("button");
    btn.innerHTML = "Remove Person";
    btn.onclick = removePerson
    btn.classList.add("remove-person")
    person.appendChild(btn);
  }
}

function removePerson() {
  this.parentNode.remove()

  let people = document.getElementById('people');
  personCount = people.childElementCount

  if (personCount == requiredParticipants) {
    removeRemovePerson()
  }
}

function removeRemovePerson() {
  els = document.getElementsByClassName("remove-person");
  while (els.length > 0) {
    els[0].parentNode.removeChild(els[0])
  }
}

function handleSubmit(event) {
  event.preventDefault();
  const data = new FormData(event.target);

  firstNames = data.getAll("first_name")
  lastNames = data.getAll("last_name")
  emails = data.getAll("email")
  count = $('.repeatable').length

  let participants = new Array()

  for (let i = 0; i < count; i++) {
    participant = new Participant(firstNames[i], lastNames[i], emails[i])
    participants.push(participant)
  }

  $.ajax({
    type: "POST",
    url: "/santa",
    data: JSON.stringify(participants),
    success: function (data) {
      $('#response').html(data)
      document.getElementById('secretSantaForm').remove()
      document.body.style.backgroundColor = "#4a934a";
    },
    error: function (data) {
      $('#response').html(data.responseText)
      document.body.style.backgroundColor = "#E57373";
    },
    dataType: "text",
    contentType: "application/json",
  });
}

function Participant(fname, lname, email) {
  this.FirstName = fname
  this.LastName = lname
  this.Email = email
}

