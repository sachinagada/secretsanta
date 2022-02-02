// html, css, and javascript copied from
// https://codepen.io/JonnyNineToes/pen/AgEax
// edited it to fit with the secret santa form


// To add javascript functionality - if javascript is enabled
document.addEventListener("DOMContentLoaded", function(){
  // insert an "add more" button after the last element
  $('#count').before('<input class="add_btn" id="add_btn" type="button" value="Add Person">');

  // add the "remove" link to each section
  $('.repeatable').prepend('<a class="remove" href="javascript:void(0);">[Remove]</a>');

  updateRemoveLinks = function () {
    // if "repeatable" element count is greater than 3...
    if ($('.repeatable').length > 3) {
      // ...show the "remove" link
      $('.repeatable').children('.remove').css({'display':'block'});
    // otherwise...
    } else {
      // don't display the "remove" link
      $('.repeatable').children('.remove').css({'display':'none'});
    }
  }

  updateRemoveLinks();

  // DESTROY!!! >.<
  // this method does all the checking necessary for deleting an element
  destroy = function () {
    // when the user clicks the "remove" link on a section...
    $('.remove').click(function(){
      // ...get that link's "repeatable" parent (container) and remove it
      $(this).parent('.repeatable').remove();
      // now, for each remaining repeatable element...
      $('.repeatable').each(function(){
        var r = this; // store "this", to avoid scope problems
        var num = $(r).index() + 1; // store the index+1 of the current "repeatable" in the collection (zero-based indexes)
        // RE-NUMBER ALL THE THINGS!!! >.<
        // change the heading to reflect the index+1 value
        $(r).children('h2').html('Person ' + num).text();
        // go through all text box elements within...
        $(r).find('input, textarea').each(function(){
          // ...change their "structure" data attributes to reflect the index+1 value of the "repeatable" element
          dattr = $(this).data('structure') + num;
           // update the "for" attribute on the parent element (label)
          $(this).attr({}).parent('label').attr('for',dattr);
        });
        //adjust the counter
        $('#count').val($('.repeatable').length);
      });
      updateRemoveLinks();
    });
  }

  // now, call the "destroy" method, so that when the page loads, this method is declared and will affect all "repeatable" elements
  destroy();

  // when the user clicks the "add more" button...
  $('.add_btn').click(function(){
    var original = $('.repeatable').last().find(':checked');
    console.log(original);
    // clone the previous element (a "repeatable" element), and insert it before the "add more" button
    $(this).prev('.repeatable').clone().insertBefore(this).html();
    // get the number of repeatable elements on the page
    var num = $('.repeatable').length;
    // again, get the previous element (a "repeatable" element), and change the header to reflect it's new index
    $(this).prev('.repeatable').children('h2').html('Person ' + num);
    // now, go through all text boxes within the last "repeatable" element...
    $('.repeatable').last().find('input, textarea').each(function(){
      // ...change their "structure" data attributes to reflect the index+1 value of the "repeatable" element
      dattr = $(this).data('structure') + num;
      // update the "for" attribute on the parent element (label)
      $(this).attr({}).parent('label').attr('for',dattr);
      // clear the input field contentqs of the new "repeatable"
      // if the type of the input is "radio"...
      if ($(this).attr('type') == 'radio') {
        // remove the checked attribute
        $(this).removeAttr('checked');
      // for all other inputs...
      } else {
        // clear the value...
        $(this).val('');
      }
      if (original.length == 1) {
          original.prop('checked',true);
      }
      //adjust the counter
      $('#count').val($('.repeatable').length);
    });
    // run the "destroy" method... I forget why... just do it, and don't gimme no lip.
    destroy();
    updateRemoveLinks();
  });

  function handleSubmit(event) {
    event.preventDefault();
    const data = new FormData(event.target);

    firstNames = data.getAll("first_name")
    lastNames = data.getAll("last_name")
    emails = data.getAll("email")
    count = data.get("count")

    let participants = new Array()

    for (let i = 0; i < count; i++) {
      participant = new Participant(firstNames[i], lastNames[i], emails[i])
      participants.push(participant)
    }

    console.log("Request successful!")

    // TODO: return response from backend to front end
    $.ajax({
      type: "POST",
      url: "/santa",
      data: JSON.stringify(participants),
      success: function(){},
      dataType: "json",
      contentType : "application/json"
    });

    // clear the data after submitting
    form.reset();
  }

  const form = document.querySelector('#myForm');
  form.addEventListener('submit', handleSubmit);
});

function Participant(fname, lname, email) {
  this.FirstName = fname
  this.LastName = lname
  this.Email = email
}
