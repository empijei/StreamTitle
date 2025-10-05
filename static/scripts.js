function addSignature(event) {
    const owner = event.target.dataset.owner;
    const user = event.target.dataset.user;
    const textarea = document.getElementById("textarea_form");
    textarea.value = `${textarea.value}\nfor ${owner} by ${user}`;
}
document.addEventListener('DOMContentLoaded', () => {
    const btn = document.getElementById("signature_button");
    // There would normally be one script per page, but let's
    // keep this simple and use an "if" instead.
    if (btn) btn.addEventListener('click', addSignature);

    const reload = document.getElementById("reload_link");
    if (reload) reload.addEventListener('click', () => {
        window.location.href = window.location.href;
    })
})