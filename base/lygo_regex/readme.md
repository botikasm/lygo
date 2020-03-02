# lygo_regex
Original code from:
https://github.com/mingrammer/commonregex

## Usage

```go

func main() {
    text := `John, please get that article on www.linkedin.com to me by 5:00PM on Jan 9th 2012. 4:00 would be ideal, actually. If you have any questions, You can reach me at (519)-236-2723x341 or get in touch with my associate at harold.smith@gmail.com`

    dateList := lygo_regex.Date(text)
    // ['Jan 9th 2012']
    timeList := lygo_regex.Time(text)
    // ['5:00PM', '4:00']
    linkList := lygo_regex.Links(text)
    // ['www.linkedin.com', 'harold.smith@gmail.com']
    phoneList := lygo_regex.PhonesWithExts(text)  
    // ['(519)-236-2723x341']
    emailList := lygo_regex.Emails(text)
    // ['harold.smith@gmail.com']
}
```

## Features

* Date
* Time
* Phone
* Phones with exts
* Link
* Email
* IPv4
* IPv6
* IP
* Ports without well-known (not known ports)
* Price
* Hex color
* Credit card
* VISA credit card
* MC credit card
* ISBN 10/13
* BTC address
* Street address
* Zip code
* Po box
* SSN
* MD5
* SHA1
* SHA256
* GUID
* MAC address
* IBAN
* Git Repository