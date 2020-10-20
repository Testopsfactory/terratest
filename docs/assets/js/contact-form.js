var params = window.location.search;
if (params === "?plan=enterprise") {
  var plan = document.getElementById("enterprise");
  plan.checked = true
} else {
  var plan = document.getElementById("pro");
  plan.checked = true
}

/* Contact form */
$(function() {
  var submitButton = $("#submit-button");
  var form = $("#contact-form");

  function submitForm() {

    clearErrors();
  
    form.find("[required]").each(function(index, el) {
      if (!$(el).val()) {
        isValid = false;
        showInputError(el);
        showFormError("Please fill in all required fields");
      }
    });
  
    function showInputError(el) {
      $(el).addClass("has-error");
    };
  
    function showFormError(message) {
      $("#error-message").html(
        '<h3 class="text-danger text-center">' + message + "</h3>"
      );
    };
  
    function clearErrors() {
      $("#error-message").html("");
      form.find("*").removeClass("has-error");
    };
  }

  submitButton.on("click", submitForm);
});