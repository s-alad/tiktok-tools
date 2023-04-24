main = document.getElementById('main');
loading = document.getElementById('loading')
loading.innerHTML = "Loading..."

function fetchImages() {
    url = "http://127.0.0.1:5000/getimagesquick"
    console.log(url)

    //fetch url and when done fetching call decode 
    fetch(url)
        .then(response => response.json())
        .then(data => decode(data))

    //remove loading from main
}

function decode(res) {
    dump = res["dump"]
    console.log(dump)
    for (i = 0; i < dump.length; i++) {
        // encode image from base64
        img = document.createElement('img')
        img.src = "data:image/png;base64," + dump[i]
        // add class square to image
        img.classList.add('square')
        //onclick the image, show the image in the center as a modal
      /*   img.onclick = function () {
            modal = document.getElementById('modal')
            modalImg = document.getElementById('modalImg')
            modal.style.display = "block"
            modalImg.src = this.src

            //when the user clicks on the modal, close it
            modal.onclick = function () {
                modal.style.display = "none"
            }
        } */

        //onclick image take to new tab
        img.onclick = function () {
            // open in new tab
            window.open(this.src)
        }
        main.appendChild(img)
    }
    loading.innerHTML = ""

}

fetchImages()