# Sprint Review Transcript

**Project:** MVP v2 - Telegram Bot and Mini App  
**Event:** Sprint Review (including customer-executed UAT)  
**Language:** English (translated from original meeting)  
**Duration:** 13:09

> **Moodle-only timecodes**
>
> - **Customer-executed UAT:** 00:03–11:54
> - **Sprint Review discussion:** 11:54–13:09

---

## Transcript

**[00:03] Speaker 1:**  
At today's Sprint Review, we planned to demonstrate the completed MVP v2 increment, including the Telegram user bot, the updated Telegram Mini App, and the admin bot. During the session, you can freely test the application, add real products to the database, and cancel any operation if necessary. Feel free to explore any available functionality.

**[00:35] Speaker 1:**  
While preparing for this meeting, I discovered several bugs in the Telegram Mini App. If you notice any additional issues during testing, please let me know. If they are already documented, I'll confirm that they are on our issue list.

**[01:11] Speaker 1:**  
I have opened the user bot. The Telegram Mini App is attached to it and can be accessed directly from the chat interface.

**[01:38] Speaker 1:**  
Could you please start screen sharing?

**[01:41] Speaker 2:**  
Sure.

**[01:42] Speaker 1:**  
Thank you.

**[01:52] Speaker 1:**  
I also remember your previous feedback about randomizing the questionnaire questions. That improvement is planned and will be implemented.

**[02:20] Speaker 1:**  
There are no changes in this section yet.

**[02:35] Speaker 1:**  
I also remember the issue with the text displayed at the bottom of this page. We'll correct that.

**[02:48] Speaker 2:**  
There is also a checkbox here.

**[03:13] Speaker 1:**  
I'd like to discuss the current sorting behavior with you. We want to remove  the sorting here, do you approve it?

**[03:22] Speaker 2:**  
Yes.

---

### Customer-Executed UAT

**[04:27] Speaker 1:**  
The address field was left empty. The application should display a validation error indicating that the address is required.

**[04:43] Speaker 1:**  
The surname field is less important, but we definitely need validation for phone numbers and email addresses. At the moment, the form can be submitted without proper validation.

**[05:06] Speaker 1:**  
I'm entering Moscow as the address for testing.

**[05:11] Speaker 2:**  
I can't continue to the payment step.

**[05:14] Speaker 1:**  
That's expected. As far as I remember, the payment process hasn't been implemented yet, so there is nowhere to continue after this screen.

**[05:21] Speaker 1:**  
Phone number and email validation should be enforced everywhere to prevent incorrect data entry.

**[05:32] Speaker 1:**  
When selecting an address, the application should also clearly display the selected location.

**[05:49] Speaker 2:**  
It actually appears here.

**[06:05] Speaker 1:**  
Good.

**[06:20] Speaker 1:**  
Why didn't the location update? I changed the first address, but it still shows Podolsk.

**[06:32] Speaker 1:**  
I'll discuss this with our frontend and backend developers because it may be related to database synchronization.

**[06:39] Speaker 2:**  
It has actually updated.

**[06:42] Speaker 1:**  
You're right. Good.

**[06:47] Speaker 1:**  
Let's also test cities like Krasnodar or Vladivostok.

**[06:53] Speaker 2:**  
Vladivostok isn't available.

**[06:57] Speaker 1:**  
Krasnodar isn't available either.

**[07:15] Speaker 1:**  
Another known issue appears when more than two products are added to the cart—the layout breaks. We'll fix that as well.

**[07:36] Speaker 2:**  
Let's test removing products from the cart.

**[07:45] Speaker 1:**  
Currently, if only one product remains, clicking remove immediately deletes it. Instead, there should be a confirmation dialog asking whether the user really wants to remove the item.

**[07:55] Speaker 2:**  
Understood.

**[08:11] Speaker 2:**  
Yes, I understand.

---

### Admin Bot Demonstration

**[08:14] Speaker 1:**  
Those are all the updates for the user application. Now let's look at the admin bot.

**[08:28] Speaker 2:**  
What's the password?

**[08:32] Speaker 1:**  
I'll send it to your inbox.

**[08:50] Speaker 1:**  
Once you're logged in, you'll see all available administrative functions.

**[09:01] Speaker 2:**  
The password isn't working.

**[09:21] Speaker 1:**  
Really? Let me try it myself.

**[09:39] Speaker 1:**  
It looks like someone may have changed the password.

**[09:52] Speaker 1:**  
Could we wait a little while until the responsible team member replies?

**[10:05] Speaker 1:**  
In the meantime, we can continue discussing the system.

---

### Product Demonstration

**[10:24] Speaker 1:**  
Earlier I added a real product to the database. Let's check whether it's still available.

**[10:37] Speaker 1:**  
Let's look at its characteristics.

**[10:53] Speaker 1:**  
Could you complete the questionnaire once again and answer every question? That should take us to the product selection screen.

**[11:35] Speaker 1:**  
It appears the product has been removed from the database.

**[11:47] Speaker 1:**  
At this point, I have demonstrated everything I intended to show.

---

## Sprint Review Discussion

**[11:54] Speaker 1:**  
Do you have any questions or anything you'd like to clarify?

**[12:04] Speaker 1:**  
I wanted to demonstrate a real product with its image, but apparently it is no longer present in the database. We'll show it again later after restoring the data.

**[12:21] Speaker 1:**  
The admin bot demonstration will also need to be completed later once access is restored.

**[12:25] Speaker 1:**  
Overall, the current implementation looks stable. The remaining issues are mostly validation improvements, UI fixes, and restoring missing data.

**[12:30] Speaker 1:**  
Regarding the admin functionality, if necessary we'll arrange another meeting on Monday or Tuesday.

**[12:37] Speaker 1:**  
You can also send me a message over the weekend if you have additional feedback.

**[12:40] Speaker 2:**  
I'll review everything in more detail and send any further comments if necessary.

**[12:42] Speaker 1:**  
Great.

**[12:45] Speaker 1:**  
As mentioned earlier, you can continue testing all available functionality. You can add real products, cancel operations at any time, and perform as much testing as needed.

**[13:02] Speaker 1:**  
That's everything from us today.

**[13:04] Speaker 1:**  
Thank you for attending the Sprint Review.

**[13:06] Speaker 2:**  
Thank you. Goodbye.

**[13:09] End of meeting**
