<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Dynamic Div Example</title>
<!-- in the template: -->
    <style>
        .plain {
            /* Default plain formatting, no additional styles */
        }

        .rounded-box {
            border: 1px solid #8B4513;
            border-radius: 15px;
            padding: 20px;
            width: 300px;
            position: relative;
            background-color: #D2B48C;
        }

        .square-box {
            border: 1px solid #800000;
            padding: 20px;
            width: 300px;
            position: relative;
            background-color: #800000;
            color: white;
        }

        .rounded-box img,
        .square-box img {
            position: absolute;
            top: 10px;
            right: 10px;
            border-radius: 50%;
            width: 50px;
            height: 50px;
        }
    </style>
</head>
<body>
<!-- what's actually in the .md: -->
    <div>some text</div>
    <div t="t2">text in box</div>
    <div t="t3">more text</div>

<!-- in the template: -->
    <script>
        document.addEventListener("DOMContentLoaded", () => {
            const divs = document.querySelectorAll("div");

            divs.forEach(div => {
                const type = div.getAttribute("t");

                if (type === "t2") {
                    div.classList.add("rounded-box");
                    const img = document.createElement("img");
                    img.src = "https://cdn.xeiaso.net/sticker/mara/hacker/128";
                    img.alt = "Img1";
                    div.appendChild(img);
                } else if (type === "t3") {
                    div.classList.add("square-box");
                    const img = document.createElement("img");
                    img.src = "https://cdn.xeiaso.net/sticker/aoi/wut/128";
                    img.alt = "Img2";
                    div.appendChild(img);
                } else {
                    div.classList.add("plain");
                }
            });
        });
    </script>
</body>
