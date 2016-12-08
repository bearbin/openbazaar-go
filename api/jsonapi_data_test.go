package api

const notFoundJSON = `{"success": false,"reason": "Not Found"}`

//
// Settings
//

const settingsJSON = `{
    "paymentDataInQR": true,
    "showNotifications": true,
    "showNsfw": true,
    "shippingAddresses": [{
        "name": "Seymour Butts",
        "company": "Globex Corporation",
        "addressLineOne": "31 Spooner Street",
        "addressLineTwo": "Apt. 124",
        "city": "Quahog",
        "state": "RI",
        "country": "UNITED_STATES",
        "postalCode": "",
        "addressNotes": "Leave package at back door"
    }],
    "localCurrency": "USD",
    "country": "UNITED_STATES",
    "language": "English",
    "termsAndConditions": "By purchasing this item you agree to the following...",
    "refundPolicy": "All sales are final.",
    "blockedNodes": ["QmecpJrN9RJ7smyYByQdZUy5mF6aapgCfKLKRmDtycv9aG", "QmamudHQGtztShX7Nc9HcczehdpGGWpFBWu2JvKWcpELxr", "QmPDLS7TV9Q3gtxRXQVqrm2RpEtz1Mq6u2YGeuEJWCqu6B"],
    "storeModerators": ["QmNedYJ6WmLhacAL2ozxb4k33Gxd9wmKB7HyoxZCwXid1e", "QmQdi7EaJUmuRUtSaCPkijw5cptFfNcX2EPvMyQwR117Y2"],
    "smtpSettings": {
        "notifications": true,
        "serverAddress": "smtp.urbanart.com:465",
        "username": "urbanart",
        "password": "letmein",
        "senderEmail": "notifications@urbanart.com",
        "recipientEmail": "Dave@gmail.com"
    }
}`

const settingsUpdateJSON = `{
    "paymentDataInQR": false,
    "showNotifications": true,
    "showNsfw": false,
    "shippingAddresses": [{
        "name": "I.C. Wiener",
        "company": "Globex Corporation",
        "addressLineOne": "31 Spooner Street",
        "addressLineTwo": "Apt. 124",
        "city": "Quahog",
        "state": "RI",
        "country": "UNITED_STATES",
        "postalCode": "",
        "addressNotes": "Leave package at front door"
    }],
    "localCurrency": "BTC",
    "country": "UNITED_STATES",
    "language": "English",
    "termsAndConditions": "By purchasing this item you agree to the following...",
    "refundPolicy": "All sales are final.",
    "blockedNodes": ["QmecpJrN9RJ7smyYByQdZUy5mF6aapgCfKLKRmDtycv9aG", "QmamudHQGtztShX7Nc9HcczehdpGGWpFBWu2JvKWcpELxr", "QmPDLS7TV9Q3gtxRXQVqrm2RpEtz1Mq6u2YGeuEJWCqu6B"],
    "storeModerators": ["QmNedYJ6WmLhacAL2ozxb4k33Gxd9wmKB7HyoxZCwXid1e", "QmQdi7EaJUmuRUtSaCPkijw5cptFfNcX2EPvMyQwR117Y2"],
    "smtpSettings": {
        "notifications": true,
        "serverAddress": "smtp.urbanart.com:465",
        "username": "urbanart",
        "password": "letmein",
        "senderEmail": "notifications@urbanart.com",
        "recipientEmail": "Dave@gmail.com"
    }
}`

const settingsMalformedJSON = `{
    /"paymentDataInQR": false,
}`

const settingsMalformedJSONErr = `{
    "success": false,
    "reason": "invalid character '/' looking for beginning of object key string"
}`

const settingsAlreadyExistsJSON = `{
    "success": false,
    "reason": "Settings is already set. Use PUT."
}`

//
// Profile
//

const profileJSON = `{
    "handle": "satoshi",
    "name": "Satoshi Nakamoto",
    "location": "Japan",
    "about": "Bitcoin's Creator",
    "shortDescription": "I make money",
    "website": "bitcoin.org",
    "email": "satoshi@gmx.com",
    "phoneNumber": "5551234567",
    "avgRating": 1,
    "numRatings": 21000000,
    "nsfw": true,
    "vendor": true,
    "moderator": true,
    "primaryColor": "#000",
    "secondaryColor": "#FFD700",
    "textColor": "#fff",
    "followerCount": 1,
    "followingCount": 2,
    "listingCount": 3,
    "lastModified": 1292110756
}`

const profileUpdateJSON = `{
    "handle": "satoshi",
    "name": "Craig Wright",
    "location": "Austrailia"
}`

const profileUpdatedJSON = `{
    "handle": "satoshi",
    "name": "Craig Wright",
    "location": "Austrailia",
    "about": "",
    "shortDescription": "",
    "website": "",
    "email": "",
    "phoneNumber": "",
    "avgRating": 0,
    "numRatings": 0,
    "nsfw": false,
    "vendor": false,
    "moderator": false,
    "primaryColor": "",
    "secondaryColor": "",
    "textColor": "",
    "followerCount": 0,
    "followingCount": 0,
    "listingCount": 0,
    "lastModified": 0
}`

const profileAlreadyExistsJSON = `{
    "success": false,
    "reason": "Profile already exists. Use PUT."
}`
