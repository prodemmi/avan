> fix iransans font for dashboard
> swagger in dev page should be force ltr for each languages and fix colors and backgrounds in dark mode too
> API Security Policy settings in json inputs are very bad! make them text inputs
> remove logs directory in settings
> make API keys per providers (not workspaces)
> fix api key copy. also show it once in a dialog after creation. it is so secure value.
> remove multi language, i want persian only
> fix all date range filters. there is one in analytics page. their place and align are not good. also i want a real persian datetime picker for them.
> buttons and elements are too big. i want more beautifull design system.
> remove message templates from system. i think it is in business sides. this system only message sender.
> i want a pro provider load balancing. currently we have it but i want to able to set wight and provider fallbacks. for example i have 3 sms providers. when i try to send a message with this provider and it got error (use retry in api for example got 3 retry errors), it should resend with fallback sms provider and it should log or saves in db and show in the table that this message got x times error (with save errors) and used y fallback providers.it was an example of sms provider only. consider any providers. in providers page, write two tabs, one providers and other is load balancing, in load balancing page show a graph of loadbalancings and enable to set or edit load balancing system here too.so i want more pro system. currently you show an wight in table and a order in provider edit, it is not pro.
> encrypy providers secure fields in database
> show workspaces in header with select box instead of sidebar
> remove "زمان و تاریخ" in settings and use default
> remove this api,/api/messages/{id}/status
> combine /send with /quick-send in one /send api, move quick_send boolean in body.
> /send should support multi providers. for example emails need title, description, attachments, recipient and ...
sms only needs recipient and text
telegram needs description, recipient id and message keyboards (optional)
> return all data of the message in GET /api/messages/{id}
> explain all type of providers /send in the swagger and say how to do that.
> add and show all table detail actions on a left drawer
> for confirms like delete, show confirm html dialog
> make dashboard welcome card smaller
