function addSignature(uname) {
    const textarea = document.getElementById("textarea_form");
    textarea.value = textarea.value + "\n" + uname;
}