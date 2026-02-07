# Feature: User Photo Promotion to Catalog

## Overview
Allow admins to see and promote user-uploaded inventory photos to become the catalog image for gear items. When a user-generated photo is promoted, the catalog item must be approved by a moderation admin before it becomes publicly visible.

## User Story
**As an admin**, I want to see photos that users have uploaded for their inventory items when moderating the gear catalog, so that I can promote high-quality user photos to become the official catalog image instead of finding/uploading images manually.

**As an admin**, I want catalog items with user-generated photos to require approval before being publicly visible, so that inappropriate or low-quality images don't appear in the public catalog.

## Acceptance Criteria

### Photo Promotion
- [ ] Gear moderation page shows a "User Photos" section for each catalog item
- [ ] Display thumbnail grid of all user-uploaded photos for inventory items linked to that catalog entry
- [ ] Each thumbnail shows which user uploaded it (optionally)
- [ ] Admin can click a user photo to preview it full-size
- [ ] Admin can click "Use as Catalog Image" to promote the user photo
- [ ] User's original inventory item keeps their photo (don't remove it)
- [ ] Moderation queue prioritizes items that have user photos available

### Approval Workflow
- [ ] When a user-generated photo is promoted to a catalog item, the item is marked as "pending approval"
- [ ] Catalog items with pending approval do NOT appear in public catalog search/browse
- [ ] Pending approval items still appear in the admin moderation queue with a "Pending" badge
- [ ] Admin can approve or reject the catalog item with the user-generated image
- [ ] On approval, the catalog item becomes publicly visible
- [ ] On rejection, the user-generated image is removed and item returns to previous state (or stays hidden if new)
- [ ] Admin can optionally add a rejection reason for internal tracking

## Technical Design

### Database Changes
1. Add `image_approval_status` column to `gear_catalog` table:
   ```sql
   ALTER TABLE gear_catalog ADD COLUMN image_approval_status VARCHAR(20) DEFAULT 'approved';
   -- Values: 'approved', 'pending', 'rejected'
   ```
2. Add `image_source` column to track where image came from:
   ```sql
   ALTER TABLE gear_catalog ADD COLUMN image_source VARCHAR(20) DEFAULT 'admin';
   -- Values: 'admin', 'user'
   ```

### Backend Changes
1. **New endpoint**: `GET /api/admin/gear-catalog/{id}/user-images`
   - Returns list of inventory items with images linked to this catalog ID
   - Includes image URL and user info

2. **Database query for user images**:
   ```sql
   SELECT i.id, i.image_url, i.user_id, u.display_name
   FROM inventory_items i
   JOIN users u ON i.user_id = u.id
   WHERE i.catalog_id = $1 AND i.image_url IS NOT NULL
   ```

3. **Promote endpoint**: `POST /api/admin/gear-catalog/{id}/promote-user-image`
   - Takes inventory item ID
   - Copies the image to catalog (either as URL or fetches/stores as binary)
   - Sets `image_source = 'user'` and `image_approval_status = 'pending'`

4. **Approval endpoint**: `POST /api/admin/gear-catalog/{id}/approve-image`
   - Sets `image_approval_status = 'approved'`

5. **Rejection endpoint**: `POST /api/admin/gear-catalog/{id}/reject-image`
   - Clears the image and sets `image_approval_status = 'approved'` (back to normal state)
   - Or optionally sets to 'rejected' for audit trail

6. **Update public catalog queries**:
   - Add `WHERE image_approval_status = 'approved'` to all public-facing catalog endpoints
   - Admin endpoints should still show all items regardless of status

### Frontend Changes
1. Add "User Photos" section to `AdminGearModeration.tsx`
2. Show grid of user-uploaded images when editing a catalog item
3. Add "Use as Catalog Image" button on each user photo
4. Show count badge on items that have user photos available
5. Add "Pending Approval" badge/filter in moderation list
6. Add "Approve" and "Reject" buttons for pending items
7. Show warning when promoting user image: "This item will require approval before appearing in public catalog"

## Dependencies
- Requires user photo upload feature to be implemented first
- Image moderation on user uploads should run before photos appear here

## Priority
Low - Nice to have for catalog curation efficiency

---

## Implementation Prompt

Use this prompt to implement the feature:

```
Implement the "User Photo Promotion to Catalog" feature for FlyingForge gear moderation with approval workflow.

Requirements:

### Database Migration
1. Add migration to gear_catalog table:
   - Add column `image_approval_status VARCHAR(20) DEFAULT 'approved'` (values: 'approved', 'pending', 'rejected')
   - Add column `image_source VARCHAR(20) DEFAULT 'admin'` (values: 'admin', 'user')

### Backend Endpoints
2. Add GET /api/admin/gear-catalog/{id}/user-images that returns all inventory items with images linked to a specific catalog ID. Include the image URL, inventory item ID, user ID, and user display name.

3. Add POST /api/admin/gear-catalog/{id}/promote-user-image that accepts an inventory item ID and:
   - Copies that item's image to become the catalog image (store as image_data binary)
   - Sets image_source = 'user'
   - Sets image_approval_status = 'pending'

4. Add POST /api/admin/gear-catalog/{id}/approve-image that sets image_approval_status = 'approved'

5. Add POST /api/admin/gear-catalog/{id}/reject-image that:
   - Clears the image_data and image_url
   - Sets image_source = 'admin' and image_approval_status = 'approved'

6. Update ALL public-facing catalog endpoints (search, list, get) to filter by image_approval_status = 'approved'. Admin endpoints should show all items.

### Frontend Changes
7. In AdminGearModeration.tsx, when editing a catalog item, add a "User Photos" section below the current image upload area:
   - Fetch user images for the current catalog item
   - Display them in a thumbnail grid
   - Each thumbnail should have a "Use as Catalog Image" button
   - Show "No user photos available" if none exist
   - Show confirmation dialog warning that item will need approval

8. Add visual indicators in the moderation list:
   - Badge showing count of user photos available per item
   - "Pending Approval" badge for items with pending status
   - Filter option to show only pending items

9. Add approval workflow UI:
   - "Approve" button (green) for pending items
   - "Reject" button (red) for pending items
   - After approval, item becomes publicly visible
   - After rejection, image is cleared and item returns to normal state

10. After promoting a user image, the "User Photos" section should still show the original user photos (don't remove them from inventory).

The backend code is in server/internal/httpapi/admin_api.go and server/internal/database/gear_catalog_store.go.
The frontend moderation component is in web/src/components/AdminGearModeration.tsx.
Database migrations go in server/internal/database/db.go.
```
